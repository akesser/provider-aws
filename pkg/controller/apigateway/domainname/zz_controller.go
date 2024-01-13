/*
Copyright 2021 The Crossplane Authors.

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

// Code generated by ack-generate. DO NOT EDIT.

package domainname

import (
	"context"

	svcapi "github.com/aws/aws-sdk-go/service/apigateway"
	svcsdk "github.com/aws/aws-sdk-go/service/apigateway"
	svcsdkapi "github.com/aws/aws-sdk-go/service/apigateway/apigatewayiface"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	cpresource "github.com/crossplane/crossplane-runtime/pkg/resource"

	svcapitypes "github.com/crossplane-contrib/provider-aws/apis/apigateway/v1alpha1"
	connectaws "github.com/crossplane-contrib/provider-aws/pkg/utils/connect/aws"
	errorutils "github.com/crossplane-contrib/provider-aws/pkg/utils/errors"
)

const (
	errUnexpectedObject = "managed resource is not an DomainName resource"

	errCreateSession = "cannot create a new session"
	errCreate        = "cannot create DomainName in AWS"
	errUpdate        = "cannot update DomainName in AWS"
	errDescribe      = "failed to describe DomainName"
	errDelete        = "failed to delete DomainName"
)

type connector struct {
	kube client.Client
	opts []option
}

func (c *connector) Connect(ctx context.Context, mg cpresource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*svcapitypes.DomainName)
	if !ok {
		return nil, errors.New(errUnexpectedObject)
	}
	sess, err := connectaws.GetConfigV1(ctx, c.kube, mg, cr.Spec.ForProvider.Region)
	if err != nil {
		return nil, errors.Wrap(err, errCreateSession)
	}
	return newExternal(c.kube, svcapi.New(sess), c.opts), nil
}

func (e *external) Observe(ctx context.Context, mg cpresource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*svcapitypes.DomainName)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}
	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}
	input := GenerateGetDomainNameInput(cr)
	if err := e.preObserve(ctx, cr, input); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "pre-observe failed")
	}
	resp, err := e.client.GetDomainNameWithContext(ctx, input)
	if err != nil {
		return managed.ExternalObservation{ResourceExists: false}, errorutils.Wrap(cpresource.Ignore(IsNotFound, err), errDescribe)
	}
	currentSpec := cr.Spec.ForProvider.DeepCopy()
	if err := e.lateInitialize(&cr.Spec.ForProvider, resp); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "late-init failed")
	}
	GenerateDomainName(resp).Status.AtProvider.DeepCopyInto(&cr.Status.AtProvider)
	upToDate := true
	diff := ""
	if !meta.WasDeleted(cr) { // There is no need to run isUpToDate if the resource is deleted
		upToDate, diff, err = e.isUpToDate(ctx, cr, resp)
		if err != nil {
			return managed.ExternalObservation{}, errors.Wrap(err, "isUpToDate check failed")
		}
	}
	return e.postObserve(ctx, cr, resp, managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        upToDate,
		Diff:                    diff,
		ResourceLateInitialized: !cmp.Equal(&cr.Spec.ForProvider, currentSpec),
	}, nil)
}

func (e *external) Create(ctx context.Context, mg cpresource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*svcapitypes.DomainName)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedObject)
	}
	cr.Status.SetConditions(xpv1.Creating())
	input := GenerateCreateDomainNameInput(cr)
	if err := e.preCreate(ctx, cr, input); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "pre-create failed")
	}
	resp, err := e.client.CreateDomainNameWithContext(ctx, input)
	if err != nil {
		return managed.ExternalCreation{}, errorutils.Wrap(err, errCreate)
	}

	if resp.CertificateArn != nil {
		cr.Spec.ForProvider.CertificateARN = resp.CertificateArn
	} else {
		cr.Spec.ForProvider.CertificateARN = nil
	}
	if resp.CertificateName != nil {
		cr.Spec.ForProvider.CertificateName = resp.CertificateName
	} else {
		cr.Spec.ForProvider.CertificateName = nil
	}
	if resp.CertificateUploadDate != nil {
		cr.Status.AtProvider.CertificateUploadDate = &metav1.Time{*resp.CertificateUploadDate}
	} else {
		cr.Status.AtProvider.CertificateUploadDate = nil
	}
	if resp.DistributionDomainName != nil {
		cr.Status.AtProvider.DistributionDomainName = resp.DistributionDomainName
	} else {
		cr.Status.AtProvider.DistributionDomainName = nil
	}
	if resp.DistributionHostedZoneId != nil {
		cr.Status.AtProvider.DistributionHostedZoneID = resp.DistributionHostedZoneId
	} else {
		cr.Status.AtProvider.DistributionHostedZoneID = nil
	}
	if resp.DomainName != nil {
		cr.Spec.ForProvider.DomainName = resp.DomainName
	} else {
		cr.Spec.ForProvider.DomainName = nil
	}
	if resp.DomainNameStatus != nil {
		cr.Status.AtProvider.DomainNameStatus = resp.DomainNameStatus
	} else {
		cr.Status.AtProvider.DomainNameStatus = nil
	}
	if resp.DomainNameStatusMessage != nil {
		cr.Status.AtProvider.DomainNameStatusMessage = resp.DomainNameStatusMessage
	} else {
		cr.Status.AtProvider.DomainNameStatusMessage = nil
	}
	if resp.EndpointConfiguration != nil {
		f8 := &svcapitypes.EndpointConfiguration{}
		if resp.EndpointConfiguration.Types != nil {
			f8f0 := []*string{}
			for _, f8f0iter := range resp.EndpointConfiguration.Types {
				var f8f0elem string
				f8f0elem = *f8f0iter
				f8f0 = append(f8f0, &f8f0elem)
			}
			f8.Types = f8f0
		}
		if resp.EndpointConfiguration.VpcEndpointIds != nil {
			f8f1 := []*string{}
			for _, f8f1iter := range resp.EndpointConfiguration.VpcEndpointIds {
				var f8f1elem string
				f8f1elem = *f8f1iter
				f8f1 = append(f8f1, &f8f1elem)
			}
			f8.VPCEndpointIDs = f8f1
		}
		cr.Spec.ForProvider.EndpointConfiguration = f8
	} else {
		cr.Spec.ForProvider.EndpointConfiguration = nil
	}
	if resp.MutualTlsAuthentication != nil {
		f9 := &svcapitypes.MutualTLSAuthenticationInput{}
		if resp.MutualTlsAuthentication.TruststoreUri != nil {
			f9.TruststoreURI = resp.MutualTlsAuthentication.TruststoreUri
		}
		if resp.MutualTlsAuthentication.TruststoreVersion != nil {
			f9.TruststoreVersion = resp.MutualTlsAuthentication.TruststoreVersion
		}
		cr.Spec.ForProvider.MutualTLSAuthentication = f9
	} else {
		cr.Spec.ForProvider.MutualTLSAuthentication = nil
	}
	if resp.OwnershipVerificationCertificateArn != nil {
		cr.Spec.ForProvider.OwnershipVerificationCertificateARN = resp.OwnershipVerificationCertificateArn
	} else {
		cr.Spec.ForProvider.OwnershipVerificationCertificateARN = nil
	}
	if resp.RegionalCertificateArn != nil {
		cr.Spec.ForProvider.RegionalCertificateARN = resp.RegionalCertificateArn
	} else {
		cr.Spec.ForProvider.RegionalCertificateARN = nil
	}
	if resp.RegionalCertificateName != nil {
		cr.Spec.ForProvider.RegionalCertificateName = resp.RegionalCertificateName
	} else {
		cr.Spec.ForProvider.RegionalCertificateName = nil
	}
	if resp.RegionalDomainName != nil {
		cr.Status.AtProvider.RegionalDomainName = resp.RegionalDomainName
	} else {
		cr.Status.AtProvider.RegionalDomainName = nil
	}
	if resp.RegionalHostedZoneId != nil {
		cr.Status.AtProvider.RegionalHostedZoneID = resp.RegionalHostedZoneId
	} else {
		cr.Status.AtProvider.RegionalHostedZoneID = nil
	}
	if resp.SecurityPolicy != nil {
		cr.Spec.ForProvider.SecurityPolicy = resp.SecurityPolicy
	} else {
		cr.Spec.ForProvider.SecurityPolicy = nil
	}
	if resp.Tags != nil {
		f16 := map[string]*string{}
		for f16key, f16valiter := range resp.Tags {
			var f16val string
			f16val = *f16valiter
			f16[f16key] = &f16val
		}
		cr.Spec.ForProvider.Tags = f16
	} else {
		cr.Spec.ForProvider.Tags = nil
	}

	return e.postCreate(ctx, cr, resp, managed.ExternalCreation{}, err)
}

func (e *external) Update(ctx context.Context, mg cpresource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*svcapitypes.DomainName)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errUnexpectedObject)
	}
	input := GenerateUpdateDomainNameInput(cr)
	if err := e.preUpdate(ctx, cr, input); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, "pre-update failed")
	}
	resp, err := e.client.UpdateDomainNameWithContext(ctx, input)
	return e.postUpdate(ctx, cr, resp, managed.ExternalUpdate{}, errorutils.Wrap(err, errUpdate))
}

func (e *external) Delete(ctx context.Context, mg cpresource.Managed) error {
	cr, ok := mg.(*svcapitypes.DomainName)
	if !ok {
		return errors.New(errUnexpectedObject)
	}
	cr.Status.SetConditions(xpv1.Deleting())
	input := GenerateDeleteDomainNameInput(cr)
	ignore, err := e.preDelete(ctx, cr, input)
	if err != nil {
		return errors.Wrap(err, "pre-delete failed")
	}
	if ignore {
		return nil
	}
	resp, err := e.client.DeleteDomainNameWithContext(ctx, input)
	return e.postDelete(ctx, cr, resp, errorutils.Wrap(cpresource.Ignore(IsNotFound, err), errDelete))
}

type option func(*external)

func newExternal(kube client.Client, client svcsdkapi.APIGatewayAPI, opts []option) *external {
	e := &external{
		kube:           kube,
		client:         client,
		preObserve:     nopPreObserve,
		postObserve:    nopPostObserve,
		lateInitialize: nopLateInitialize,
		isUpToDate:     alwaysUpToDate,
		preCreate:      nopPreCreate,
		postCreate:     nopPostCreate,
		preDelete:      nopPreDelete,
		postDelete:     nopPostDelete,
		preUpdate:      nopPreUpdate,
		postUpdate:     nopPostUpdate,
	}
	for _, f := range opts {
		f(e)
	}
	return e
}

type external struct {
	kube           client.Client
	client         svcsdkapi.APIGatewayAPI
	preObserve     func(context.Context, *svcapitypes.DomainName, *svcsdk.GetDomainNameInput) error
	postObserve    func(context.Context, *svcapitypes.DomainName, *svcsdk.DomainName, managed.ExternalObservation, error) (managed.ExternalObservation, error)
	lateInitialize func(*svcapitypes.DomainNameParameters, *svcsdk.DomainName) error
	isUpToDate     func(context.Context, *svcapitypes.DomainName, *svcsdk.DomainName) (bool, string, error)
	preCreate      func(context.Context, *svcapitypes.DomainName, *svcsdk.CreateDomainNameInput) error
	postCreate     func(context.Context, *svcapitypes.DomainName, *svcsdk.DomainName, managed.ExternalCreation, error) (managed.ExternalCreation, error)
	preDelete      func(context.Context, *svcapitypes.DomainName, *svcsdk.DeleteDomainNameInput) (bool, error)
	postDelete     func(context.Context, *svcapitypes.DomainName, *svcsdk.DeleteDomainNameOutput, error) error
	preUpdate      func(context.Context, *svcapitypes.DomainName, *svcsdk.UpdateDomainNameInput) error
	postUpdate     func(context.Context, *svcapitypes.DomainName, *svcsdk.DomainName, managed.ExternalUpdate, error) (managed.ExternalUpdate, error)
}

func nopPreObserve(context.Context, *svcapitypes.DomainName, *svcsdk.GetDomainNameInput) error {
	return nil
}

func nopPostObserve(_ context.Context, _ *svcapitypes.DomainName, _ *svcsdk.DomainName, obs managed.ExternalObservation, err error) (managed.ExternalObservation, error) {
	return obs, err
}
func nopLateInitialize(*svcapitypes.DomainNameParameters, *svcsdk.DomainName) error {
	return nil
}
func alwaysUpToDate(context.Context, *svcapitypes.DomainName, *svcsdk.DomainName) (bool, string, error) {
	return true, "", nil
}

func nopPreCreate(context.Context, *svcapitypes.DomainName, *svcsdk.CreateDomainNameInput) error {
	return nil
}
func nopPostCreate(_ context.Context, _ *svcapitypes.DomainName, _ *svcsdk.DomainName, cre managed.ExternalCreation, err error) (managed.ExternalCreation, error) {
	return cre, err
}
func nopPreDelete(context.Context, *svcapitypes.DomainName, *svcsdk.DeleteDomainNameInput) (bool, error) {
	return false, nil
}
func nopPostDelete(_ context.Context, _ *svcapitypes.DomainName, _ *svcsdk.DeleteDomainNameOutput, err error) error {
	return err
}
func nopPreUpdate(context.Context, *svcapitypes.DomainName, *svcsdk.UpdateDomainNameInput) error {
	return nil
}
func nopPostUpdate(_ context.Context, _ *svcapitypes.DomainName, _ *svcsdk.DomainName, upd managed.ExternalUpdate, err error) (managed.ExternalUpdate, error) {
	return upd, err
}
