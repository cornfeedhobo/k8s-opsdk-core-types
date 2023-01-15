/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"context"
	"encoding/json"
	"net/http"

	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//+kubebuilder:webhook:path=/mutate-v1-job,mutating=true,failurePolicy=fail,sideEffects=None,groups=batch,resources=jobs,verbs=create;update,versions=v1,name=mjob.kb.io,admissionReviewVersions=v1

// log is for logging in this package.
var joblog = logf.Log.WithName("job-resource")

type JobMutator struct {
	Client  client.Client
	decoder *admission.Decoder
}

func RegisterJobMutatorWebhook(mgr manager.Manager) {
	joblog.Info("Registering JobMutator")
	m := JobMutator{
		Client: mgr.GetClient(),
	}
	mgr.GetWebhookServer().Register("/mutate-v1-job", &webhook.Admission{Handler: &m})
}

func (m *JobMutator) InjectDecoder(d *admission.Decoder) error {
	m.decoder = d
	return nil
}

func (m *JobMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	j := &batchv1.Job{}

	err := m.decoder.Decode(req, j)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if j.Annotations == nil {
		j.Annotations = map[string]string{}
	}
	j.Annotations["example-mutating-admission-webhook"] = "true"

	marshaled, err := json.Marshal(j)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaled)
}
