package v1

import (
	"context"
	"time"

	"github.com/davecgh/go-spew/spew"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Job", func() {
	It("maybe", func() {
		nn := types.NamespacedName{Namespace: "unittest", Name: "job"}

		orig := &batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      nn.Name,
				Namespace: nn.Namespace,
				Labels:    map[string]string{"app": "unittest"},
			},
			Spec: batchv1.JobSpec{
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "unittest"},
					},
					Spec: corev1.PodSpec{
						RestartPolicy: corev1.RestartPolicyOnFailure,
						Containers: []corev1.Container{
							{
								Name:  "main",
								Image: "example.com/unittest:latest",
							},
						},
					},
				},
			},
		}
		err := k8sClient.Create(context.TODO(), orig)
		Expect(err).NotTo(HaveOccurred())

		new := &batchv1.Job{}
		Eventually(func(g Gomega) {
			err := k8sClient.Get(context.TODO(), nn, new)
			Expect(err).NotTo(HaveOccurred())

			Expect(new.Annotations).ToNot(BeNil())
			spew.Dump(new.Annotations)
			value, ok := new.Annotations["example-mutating-admission-webhook"]
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("true"))
		}, time.Second*10, time.Millisecond*250).Should(Succeed())
	})
})
