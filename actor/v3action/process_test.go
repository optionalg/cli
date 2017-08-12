package v3action_test

import (
	"errors"
	"time"

	. "code.cloudfoundry.org/cli/actor/v3action"
	"code.cloudfoundry.org/cli/actor/v3action/v3actionfakes"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Process Actions", func() {
	var (
		actor                     *Actor
		fakeCloudControllerClient *v3actionfakes.FakeCloudControllerClient
	)

	BeforeEach(func() {
		fakeCloudControllerClient = new(v3actionfakes.FakeCloudControllerClient)
		actor = NewActor(fakeCloudControllerClient, nil)
	})

	Describe("Instance", func() {
		Describe("StartTime", func() {
			It("returns the time that the instance started", func() {
				instance := Instance{Uptime: 86400}
				Expect(instance.StartTime()).To(BeTemporally("~", time.Now().Add(-24*time.Hour), 10*time.Second))
			})
		})
	})

	Describe("Process", func() {
		var process Process
		BeforeEach(func() {
			process = Process{
				Instances: []Instance{
					Instance{State: "RUNNING"},
					Instance{State: "RUNNING"},
					Instance{State: "STOPPED"},
				},
			}
		})

		Describe("TotalInstanceCount", func() {
			It("returns the total number of instances", func() {
				Expect(process.TotalInstanceCount()).To(Equal(3))
			})
		})

		Describe("HealthyInstanceCount", func() {
			It("returns the total number of RUNNING instances", func() {
				Expect(process.HealthyInstanceCount()).To(Equal(2))
			})
		})
	})

	Describe("Processes", func() {
		var processes Processes

		BeforeEach(func() {
			processes = Processes{
				{
					Type: "worker",
					Instances: []Instance{
						{State: "RUNNING"},
						{State: "STOPPED"},
					},
				},
				{
					Type: "console",
					Instances: []Instance{
						{State: "RUNNING"},
					},
				},
				{
					Type: "web",
					Instances: []Instance{
						{State: "RUNNING"},
						{State: "RUNNING"},
						{State: "STOPPED"},
					},
				},
			}
		})

		Describe("Sort", func() {
			It("sorts processes with web first and then alphabetically sorted", func() {
				processes.Sort()
				Expect(processes[0].Type).To(Equal("web"))
				Expect(processes[1].Type).To(Equal("console"))
				Expect(processes[2].Type).To(Equal("worker"))
			})
		})

		Describe("Summary", func() {
			It("returns all processes and their instance count ", func() {
				Expect(processes.Summary()).To(Equal("web:2/3, console:1/1, worker:1/2"))
			})
		})
	})

	Describe("ScaleProcessByApplication", func() {
		var passedProcess ccv3.Process

		BeforeEach(func() {
			passedProcess = ccv3.Process{
				Type:       "web",
				Instances:  2,
				MemoryInMB: 100,
				DiskInMB:   200,
			}
		})

		Context("when no errors are encountered scaling the application process", func() {
			BeforeEach(func() {
				fakeCloudControllerClient.CreateApplicationProcessScaleReturns(
					ccv3.Process{
						GUID:       "some-process-guid",
						Type:       passedProcess.Type,
						Instances:  passedProcess.Instances,
						MemoryInMB: passedProcess.MemoryInMB,
						DiskInMB:   passedProcess.DiskInMB,
					},
					ccv3.Warnings{"scale-process-warning"},
					nil)
			})

			Context("when no errors are encountered getting process instances", func() {
				var (
					instance1 ccv3.Instance
					instance2 ccv3.Instance
				)

				BeforeEach(func() {
					instance1 = ccv3.Instance{
						Index:       0,
						State:       "RUNNING",
						CPU:         0.33,
						MemoryUsage: 10 * 1024 * 1024,
						DiskUsage:   20 * 1024 * 1024,
						MemoryQuota: 100 * 1024 * 1024,
						DiskQuota:   200 * 1024 * 1024,
					}
					instance2 = ccv3.Instance{
						Index:       1,
						State:       "RUNNING",
						CPU:         0.40,
						MemoryUsage: 10 * 1024 * 1024,
						DiskUsage:   20 * 1024 * 1024,
						MemoryQuota: 100 * 1024 * 1024,
						DiskQuota:   200 * 1024 * 1024,
					}
					fakeCloudControllerClient.GetProcessInstancesReturns(
						[]ccv3.Instance{instance1, instance2},
						ccv3.Warnings{"get-instances-warning"},
						nil)
				})

				It("returns the process with instance information and all warnings", func() {
					process, warnings, err := actor.ScaleProcessByApplication("some-app-guid", passedProcess)

					Expect(err).ToNot(HaveOccurred())
					Expect(warnings).To(ConsistOf("scale-process-warning", "get-instances-warning"))
					Expect(process).To(Equal(Process{
						Type:       "web",
						Instances:  []Instance{Instance(instance1), Instance(instance2)},
						MemoryInMB: 100,
						DiskInMB:   200,
					}))

					Expect(fakeCloudControllerClient.CreateApplicationProcessScaleCallCount()).To(Equal(1))
					appGUIDArg, processArg := fakeCloudControllerClient.CreateApplicationProcessScaleArgsForCall(0)
					Expect(appGUIDArg).To(Equal("some-app-guid"))
					Expect(processArg).To(Equal(passedProcess))

					Expect(fakeCloudControllerClient.GetProcessInstancesCallCount()).To(Equal(1))
					Expect(fakeCloudControllerClient.GetProcessInstancesArgsForCall(0)).To(Equal("some-process-guid"))
				})
			})

			Context("when an error is encountered getting process instances", func() {
				var expectedErr error

				BeforeEach(func() {
					expectedErr = errors.New("get instances error")
					fakeCloudControllerClient.GetProcessInstancesReturns(
						nil,
						ccv3.Warnings{"get-instances-warning"},
						expectedErr)
				})

				It("returns the error and all warnings", func() {
					_, warnings, err := actor.ScaleProcessByApplication("some-app-guid", passedProcess)
					Expect(err).To(MatchError(expectedErr))
					Expect(warnings).To(ConsistOf("scale-process-warning", "get-instances-warning"))
				})
			})
		})

		Context("when an error is encountered scaling the application process", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = errors.New("scale process error")
				fakeCloudControllerClient.CreateApplicationProcessScaleReturns(
					ccv3.Process{},
					ccv3.Warnings{"scale-process-warning"},
					expectedErr)
			})

			It("returns the error and all warnings", func() {
				_, warnings, err := actor.ScaleProcessByApplication("some-app-guid", passedProcess)
				Expect(err).To(MatchError(expectedErr))
				Expect(warnings).To(ConsistOf("scale-process-warning"))
			})
		})
	})
})
