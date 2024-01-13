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

package subnetgroup

import (
	"context"

	svcapi "github.com/aws/aws-sdk-go/service/dax"
	svcsdk "github.com/aws/aws-sdk-go/service/dax"
	svcsdkapi "github.com/aws/aws-sdk-go/service/dax/daxiface"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	cpresource "github.com/crossplane/crossplane-runtime/pkg/resource"

	svcapitypes "github.com/crossplane-contrib/provider-aws/apis/dax/v1alpha1"
	connectaws "github.com/crossplane-contrib/provider-aws/pkg/utils/connect/aws"
	errorutils "github.com/crossplane-contrib/provider-aws/pkg/utils/errors"
)

const (
	errUnexpectedObject = "managed resource is not an SubnetGroup resource"

	errCreateSession = "cannot create a new session"
	errCreate        = "cannot create SubnetGroup in AWS"
	errUpdate        = "cannot update SubnetGroup in AWS"
	errDescribe      = "failed to describe SubnetGroup"
	errDelete        = "failed to delete SubnetGroup"
)

type connector struct {
	kube client.Client
	opts []option
}

func (c *connector) Connect(ctx context.Context, mg cpresource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*svcapitypes.SubnetGroup)
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
	cr, ok := mg.(*svcapitypes.SubnetGroup)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}
	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}
	input := GenerateDescribeSubnetGroupsInput(cr)
	if err := e.preObserve(ctx, cr, input); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "pre-observe failed")
	}
	resp, err := e.client.DescribeSubnetGroupsWithContext(ctx, input)
	if err != nil {
		return managed.ExternalObservation{ResourceExists: false}, errorutils.Wrap(cpresource.Ignore(IsNotFound, err), errDescribe)
	}
	resp = e.filterList(cr, resp)
	if len(resp.SubnetGroups) == 0 {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	currentSpec := cr.Spec.ForProvider.DeepCopy()
	if err := e.lateInitialize(&cr.Spec.ForProvider, resp); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "late-init failed")
	}
	GenerateSubnetGroup(resp).Status.AtProvider.DeepCopyInto(&cr.Status.AtProvider)
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
	cr, ok := mg.(*svcapitypes.SubnetGroup)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedObject)
	}
	cr.Status.SetConditions(xpv1.Creating())
	input := GenerateCreateSubnetGroupInput(cr)
	if err := e.preCreate(ctx, cr, input); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "pre-create failed")
	}
	resp, err := e.client.CreateSubnetGroupWithContext(ctx, input)
	if err != nil {
		return managed.ExternalCreation{}, errorutils.Wrap(err, errCreate)
	}

	if resp.SubnetGroup.Description != nil {
		cr.Spec.ForProvider.Description = resp.SubnetGroup.Description
	} else {
		cr.Spec.ForProvider.Description = nil
	}
	if resp.SubnetGroup.SubnetGroupName != nil {
		cr.Status.AtProvider.SubnetGroupName = resp.SubnetGroup.SubnetGroupName
	} else {
		cr.Status.AtProvider.SubnetGroupName = nil
	}
	if resp.SubnetGroup.Subnets != nil {
		f2 := []*svcapitypes.Subnet{}
		for _, f2iter := range resp.SubnetGroup.Subnets {
			f2elem := &svcapitypes.Subnet{}
			if f2iter.SubnetAvailabilityZone != nil {
				f2elem.SubnetAvailabilityZone = f2iter.SubnetAvailabilityZone
			}
			if f2iter.SubnetIdentifier != nil {
				f2elem.SubnetIdentifier = f2iter.SubnetIdentifier
			}
			f2 = append(f2, f2elem)
		}
		cr.Status.AtProvider.Subnets = f2
	} else {
		cr.Status.AtProvider.Subnets = nil
	}
	if resp.SubnetGroup.VpcId != nil {
		cr.Status.AtProvider.VPCID = resp.SubnetGroup.VpcId
	} else {
		cr.Status.AtProvider.VPCID = nil
	}

	return e.postCreate(ctx, cr, resp, managed.ExternalCreation{}, err)
}

func (e *external) Update(ctx context.Context, mg cpresource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*svcapitypes.SubnetGroup)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errUnexpectedObject)
	}
	input := GenerateUpdateSubnetGroupInput(cr)
	if err := e.preUpdate(ctx, cr, input); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, "pre-update failed")
	}
	resp, err := e.client.UpdateSubnetGroupWithContext(ctx, input)
	return e.postUpdate(ctx, cr, resp, managed.ExternalUpdate{}, errorutils.Wrap(err, errUpdate))
}

func (e *external) Delete(ctx context.Context, mg cpresource.Managed) error {
	cr, ok := mg.(*svcapitypes.SubnetGroup)
	if !ok {
		return errors.New(errUnexpectedObject)
	}
	cr.Status.SetConditions(xpv1.Deleting())
	input := GenerateDeleteSubnetGroupInput(cr)
	ignore, err := e.preDelete(ctx, cr, input)
	if err != nil {
		return errors.Wrap(err, "pre-delete failed")
	}
	if ignore {
		return nil
	}
	resp, err := e.client.DeleteSubnetGroupWithContext(ctx, input)
	return e.postDelete(ctx, cr, resp, errorutils.Wrap(cpresource.Ignore(IsNotFound, err), errDelete))
}

type option func(*external)

func newExternal(kube client.Client, client svcsdkapi.DAXAPI, opts []option) *external {
	e := &external{
		kube:           kube,
		client:         client,
		preObserve:     nopPreObserve,
		postObserve:    nopPostObserve,
		lateInitialize: nopLateInitialize,
		isUpToDate:     alwaysUpToDate,
		filterList:     nopFilterList,
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
	client         svcsdkapi.DAXAPI
	preObserve     func(context.Context, *svcapitypes.SubnetGroup, *svcsdk.DescribeSubnetGroupsInput) error
	postObserve    func(context.Context, *svcapitypes.SubnetGroup, *svcsdk.DescribeSubnetGroupsOutput, managed.ExternalObservation, error) (managed.ExternalObservation, error)
	filterList     func(*svcapitypes.SubnetGroup, *svcsdk.DescribeSubnetGroupsOutput) *svcsdk.DescribeSubnetGroupsOutput
	lateInitialize func(*svcapitypes.SubnetGroupParameters, *svcsdk.DescribeSubnetGroupsOutput) error
	isUpToDate     func(context.Context, *svcapitypes.SubnetGroup, *svcsdk.DescribeSubnetGroupsOutput) (bool, string, error)
	preCreate      func(context.Context, *svcapitypes.SubnetGroup, *svcsdk.CreateSubnetGroupInput) error
	postCreate     func(context.Context, *svcapitypes.SubnetGroup, *svcsdk.CreateSubnetGroupOutput, managed.ExternalCreation, error) (managed.ExternalCreation, error)
	preDelete      func(context.Context, *svcapitypes.SubnetGroup, *svcsdk.DeleteSubnetGroupInput) (bool, error)
	postDelete     func(context.Context, *svcapitypes.SubnetGroup, *svcsdk.DeleteSubnetGroupOutput, error) error
	preUpdate      func(context.Context, *svcapitypes.SubnetGroup, *svcsdk.UpdateSubnetGroupInput) error
	postUpdate     func(context.Context, *svcapitypes.SubnetGroup, *svcsdk.UpdateSubnetGroupOutput, managed.ExternalUpdate, error) (managed.ExternalUpdate, error)
}

func nopPreObserve(context.Context, *svcapitypes.SubnetGroup, *svcsdk.DescribeSubnetGroupsInput) error {
	return nil
}
func nopPostObserve(_ context.Context, _ *svcapitypes.SubnetGroup, _ *svcsdk.DescribeSubnetGroupsOutput, obs managed.ExternalObservation, err error) (managed.ExternalObservation, error) {
	return obs, err
}
func nopFilterList(_ *svcapitypes.SubnetGroup, list *svcsdk.DescribeSubnetGroupsOutput) *svcsdk.DescribeSubnetGroupsOutput {
	return list
}

func nopLateInitialize(*svcapitypes.SubnetGroupParameters, *svcsdk.DescribeSubnetGroupsOutput) error {
	return nil
}
func alwaysUpToDate(context.Context, *svcapitypes.SubnetGroup, *svcsdk.DescribeSubnetGroupsOutput) (bool, string, error) {
	return true, "", nil
}

func nopPreCreate(context.Context, *svcapitypes.SubnetGroup, *svcsdk.CreateSubnetGroupInput) error {
	return nil
}
func nopPostCreate(_ context.Context, _ *svcapitypes.SubnetGroup, _ *svcsdk.CreateSubnetGroupOutput, cre managed.ExternalCreation, err error) (managed.ExternalCreation, error) {
	return cre, err
}
func nopPreDelete(context.Context, *svcapitypes.SubnetGroup, *svcsdk.DeleteSubnetGroupInput) (bool, error) {
	return false, nil
}
func nopPostDelete(_ context.Context, _ *svcapitypes.SubnetGroup, _ *svcsdk.DeleteSubnetGroupOutput, err error) error {
	return err
}
func nopPreUpdate(context.Context, *svcapitypes.SubnetGroup, *svcsdk.UpdateSubnetGroupInput) error {
	return nil
}
func nopPostUpdate(_ context.Context, _ *svcapitypes.SubnetGroup, _ *svcsdk.UpdateSubnetGroupOutput, upd managed.ExternalUpdate, err error) (managed.ExternalUpdate, error) {
	return upd, err
}
