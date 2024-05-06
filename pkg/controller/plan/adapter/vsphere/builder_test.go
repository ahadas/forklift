package vsphere

import (
	v1beta1 "github.com/konveyor/forklift-controller/pkg/apis/forklift/v1beta1"
	plancontext "github.com/konveyor/forklift-controller/pkg/controller/plan/context"
	"github.com/konveyor/forklift-controller/pkg/controller/provider/model/vsphere"
	model "github.com/konveyor/forklift-controller/pkg/controller/provider/web/vsphere"
	"github.com/konveyor/forklift-controller/pkg/lib/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vmware/govmomi/vim25/types"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var builderLog = logging.WithName("vsphere-builder-test")

const ManualOrigin = string(types.NetIpConfigInfoIpAddressOriginManual)

var _ = Describe("vSphere builder", func() {
	builder := createBuilder()
	DescribeTable("should", func(vm *model.VM, outputMap string) {
		Expect(builder.mapMacStaticIps(vm)).Should(Equal(outputMap))
	},
		Entry("no static ips", &model.VM{}, ""),
		Entry("single static ip", &model.VM{GuestNetworks: []vsphere.GuestNetwork{{MAC: "00:50:56:83:25:47", IP: "172.29.3.193", Origin: ManualOrigin}}}, "00:50:56:83:25:47:ip:172.29.3.193"),
		Entry("multiple static ips", &model.VM{GuestNetworks: []vsphere.GuestNetwork{
			{
				MAC:    "00:50:56:83:25:47",
				IP:     "172.29.3.193",
				Origin: ManualOrigin,
			},
			{
				MAC:    "00:50:56:83:25:47",
				IP:     "fe80::5da:b7a5:e0a2:a097",
				Origin: ManualOrigin,
			},
		},
		}, "00:50:56:83:25:47:ip:172.29.3.193_00:50:56:83:25:47:ip:fe80::5da:b7a5:e0a2:a097"),
		Entry("non-static ip", &model.VM{GuestNetworks: []vsphere.GuestNetwork{{MAC: "00:50:56:83:25:47", IP: "172.29.3.193", Origin: string(types.NetIpConfigInfoIpAddressOriginDhcp)}}}, ""),
	)
})

func createBuilder(objs ...runtime.Object) *Builder {
	scheme := runtime.NewScheme()
	_ = v1.AddToScheme(scheme)
	v1beta1.SchemeBuilder.AddToScheme(scheme)
	client := fake.NewClientBuilder().
		WithScheme(scheme).
		WithRuntimeObjects(objs...).
		Build()
	return &Builder{
		Context: &plancontext.Context{
			Destination: plancontext.Destination{
				Client: client,
			},
			Plan: createPlan(),
			Log:  builderLog,

			// To make sure r.Scheme is not nil
			Client: client,
		},
	}
}
