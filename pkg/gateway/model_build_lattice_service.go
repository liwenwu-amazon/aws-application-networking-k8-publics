package gateway

import (
	"context"
	"errors"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/aws/aws-application-networking-k8s/pkg/utils/gwlog"

	"github.com/aws/aws-application-networking-k8s/pkg/config"
	"github.com/aws/aws-application-networking-k8s/pkg/k8s"
	"github.com/aws/aws-application-networking-k8s/pkg/model/core"
	latticemodel "github.com/aws/aws-application-networking-k8s/pkg/model/lattice"

	lattice_aws "github.com/aws/aws-application-networking-k8s/pkg/aws"
	"github.com/aws/aws-application-networking-k8s/pkg/latticestore"
)

const (
	resourceIDLatticeService = "LatticeService"
)

type LatticeServiceBuilder interface {
	Build(ctx context.Context, httpRoute core.Route) (core.Stack, *latticemodel.Service, error)
}

type LatticeServiceModelBuilder struct {
	log         gwlog.Logger
	client      client.Client
	defaultTags map[string]string
	datastore   *latticestore.LatticeDataStore
	cloud       lattice_aws.Cloud
}

func NewLatticeServiceBuilder(
	log gwlog.Logger,
	client client.Client,
	datastore *latticestore.LatticeDataStore,
	cloud lattice_aws.Cloud,
) *LatticeServiceModelBuilder {
	return &LatticeServiceModelBuilder{
		log:       log,
		client:    client,
		datastore: datastore,
		cloud:     cloud,
	}
}

func (b *LatticeServiceModelBuilder) Build(
	ctx context.Context,
	route core.Route,
) (core.Stack, *latticemodel.Service, error) {
	stack := core.NewDefaultStack(core.StackID(k8s.NamespacedName(route.K8sObject())))

	task := &latticeServiceModelBuildTask{
		log:       b.log,
		route:     route,
		stack:     stack,
		client:    b.client,
		tgByResID: make(map[string]*latticemodel.TargetGroup),
		datastore: b.datastore,
	}

	if err := task.run(ctx); err != nil {
		return stack, task.latticeService, errors.New("LATTICE_RETRY")
	}

	return task.stack, task.latticeService, nil
}

func (t *latticeServiceModelBuildTask) run(ctx context.Context) error {
	err := t.buildModel(ctx)
	return err
}

func (t *latticeServiceModelBuildTask) buildModel(ctx context.Context) error {
	err := t.buildLatticeService(ctx)

	if err != nil {
		return fmt.Errorf("latticeServiceModelBuildTask: Failed on buildLatticeService %w", err)
	}

	_, err = t.buildTargetGroup(ctx, t.client)

	if err != nil {
		return fmt.Errorf("latticeServiceModelBuildTask: Failed on buildTargetGroup %w", err)
	}

	if !t.route.DeletionTimestamp().IsZero() {
		t.log.Infof("latticeServiceModelBuildTask: for delete ignore Targets, policy %v\n", t.route)
		return nil
	}

	err = t.buildTargets(ctx)

	if err != nil {
		t.log.Infof("latticeServiceModelBuildTask: Failed on building targets, error = %v\n ", err)
	}
	// only build listener when it is NOT delete case
	err = t.buildListener(ctx)

	if err != nil {
		return fmt.Errorf("latticeServiceModelBuildTask: Failed on building listener %w", err)
	}

	err = t.buildRules(ctx)

	if err != nil {
		return fmt.Errorf("latticeServiceModelBuildTask: Failed on building rule %w", err)
	}

	return nil
}

func (t *latticeServiceModelBuildTask) buildLatticeService(ctx context.Context) error {
	pro := "HTTP"
	protocols := []*string{&pro}
	spec := latticemodel.ServiceSpec{
		Name:      t.route.Name(),
		Namespace: t.route.Namespace(),
		Protocols: protocols,
		//ServiceNetworkNames: string(t.route.Spec().ParentRefs()[0].Name),
	}

	for _, parentRef := range t.route.Spec().ParentRefs() {
		spec.ServiceNetworkNames = append(spec.ServiceNetworkNames, string(parentRef.Name))
	}
	defaultGateway, err := config.GetClusterLocalGateway()
	if err == nil {
		spec.ServiceNetworkNames = append(spec.ServiceNetworkNames, defaultGateway)
	}

	if len(t.route.Spec().Hostnames()) > 0 {
		// The 1st hostname will be used as lattice customer-domain-name
		spec.CustomerDomainName = string(t.route.Spec().Hostnames()[0])

		t.log.Infof("Setting customer-domain-name: %v for route %v-%v",
			spec.CustomerDomainName, t.route.Name(), t.route.Namespace())
	} else {
		t.log.Infof("No custom-domain-name for route :%v-%v",
			t.route.Name(), t.route.Namespace())
		spec.CustomerDomainName = ""
	}

	if t.route.DeletionTimestamp().IsZero() {
		spec.IsDeleted = false
	} else {
		spec.IsDeleted = true
	}

	serviceResourceName := fmt.Sprintf("%s-%s", t.route.Name(), t.route.Namespace())

	t.latticeService = latticemodel.NewLatticeService(t.stack, serviceResourceName, spec)

	return nil
}

type latticeServiceModelBuildTask struct {
	log             gwlog.Logger
	route           core.Route
	client          client.Client
	latticeService  *latticemodel.Service
	tgByResID       map[string]*latticemodel.TargetGroup
	listenerByResID map[string]*latticemodel.Listener
	rulesByResID    map[string]*latticemodel.Rule
	stack           core.Stack
	datastore       *latticestore.LatticeDataStore
	cloud           lattice_aws.Cloud
}
