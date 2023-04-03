/*
Copyright 2019 The Kubernetes Authors.

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

package noderefutil

import (
	"testing"

	. "github.com/onsi/gomega"
)

const aws = "aws"
const azure = "azure"

func TestInvalidProviderID(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		err   error
	}{
		{
			name:  "empty id",
			input: "",
			err:   ErrEmptyProviderID,
		},
		{
			name:  "only empty segments",
			input: "aws:///////",
			err:   ErrInvalidProviderID,
		},
		{
			name:  "missing cloud provider",
			input: "://instance-id",
			err:   ErrInvalidProviderID,
		},
		{
			name:  "missing cloud provider and colon",
			input: "//instance-id",
			err:   ErrInvalidProviderID,
		},
		{
			name:  "missing cloud provider, colon, one leading slash",
			input: "/instance-id",
			err:   ErrInvalidProviderID,
		},
		{
			name:  "just an id",
			input: "instance-id",
			err:   ErrInvalidProviderID,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			g := NewWithT(t)

			_, err := NewProviderID(test.input)
			g.Expect(err).To(MatchError(test.err))
		})
	}
}

func TestProviderIDEquals(t *testing.T) {
	g := NewWithT(t)

	inputAWS1 := "aws:////instance-id1"
	parsedAWS1, err := NewProviderID(inputAWS1)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(parsedAWS1.String()).To(Equal(inputAWS1))

	inputAWS2 := "aws:///us-west-1/instance-id1"
	parsedAWS2, err := NewProviderID(inputAWS2)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(parsedAWS2.String()).To(Equal(inputAWS2))

	// Test for inequality
	g.Expect(parsedAWS1.Equals(parsedAWS2)).To(BeFalse())

	inputAzure1 := "azure:///subscriptions/4920076a-ba9f-11ec-8422-0242ac120002/resourceGroups/default-template/providers/Microsoft.Compute/virtualMachines/default-template-control-plane-fhrvh"
	parsedAzure1, err := NewProviderID(inputAzure1)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(parsedAzure1.String()).To(Equal(inputAzure1))

	inputAzure2 := inputAzure1
	parsedAzure2, err := NewProviderID(inputAzure2)
	g.Expect(err).NotTo(HaveOccurred())

	// Test for equality
	g.Expect(parsedAzure1.Equals(parsedAzure2)).To(BeTrue())

	// Here we ensure that two different ProviderID strings that happen to have the same 'ID' are not equal
	// We use Azure VMSS as an example, two different '0' VMs in different pools: k8s-pool1-vmss, and k8s-pool2-vmss
	inputAzureVMFromOneVMSS := "azure:///subscriptions/4920076a-ba9f-11ec-8422-0242ac120002/resourceGroups/default-template/providers/Microsoft.Compute/virtualMachineScaleSets/k8s-pool1-vmss/virtualMachines/0"
	inputAzureVMFromAnotherVMSS := "azure:///subscriptions/4920076a-ba9f-11ec-8422-0242ac120002/resourceGroups/default-template/providers/Microsoft.Compute/virtualMachineScaleSets/k8s-pool2-vmss/virtualMachines/0"
	parsedAzureVMFromOneVMSS, err := NewProviderID(inputAzureVMFromOneVMSS)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(parsedAzureVMFromOneVMSS.String()).To(Equal(inputAzureVMFromOneVMSS))

	parsedAzureVMFromAnotherVMSS, err := NewProviderID(inputAzureVMFromAnotherVMSS)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(parsedAzureVMFromAnotherVMSS.String()).To(Equal(inputAzureVMFromAnotherVMSS))

	// Test for inequality
	g.Expect(parsedAzureVMFromOneVMSS.Equals(parsedAzureVMFromAnotherVMSS)).To(BeFalse())
}
