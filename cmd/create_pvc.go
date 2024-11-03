package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	cmdcreate "k8s.io/kubectl/pkg/cmd/create"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/scheme"
)

type CreatePVCOptions struct {
	PrintFlags *genericclioptions.PrintFlags
	PrintObj   func(obj runtime.Object) error

	Name         string
	Namespace    string
	StorageClass string
	Size         string

	Client corev1client.CoreV1Interface

	DryRunStrategy cmdutil.DryRunStrategy

	genericiooptions.IOStreams
}

func NewCreatePVCOptions(ioStreams genericclioptions.IOStreams) *CreatePVCOptions {
	return &CreatePVCOptions{
		PrintFlags: genericclioptions.NewPrintFlags("created").WithTypeSetter(scheme.Scheme),
		IOStreams:  ioStreams,
	}
}

func NewCmdCreatePVC(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewCreatePVCOptions(ioStreams)

	var pvcCmd = &cobra.Command{
		Use:   "pvc NAME --storageclass <sc> ...",
		Short: "Create PVC resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Run())
		},
	}

	pvcCmd.Flags().StringVarP(&o.StorageClass, "storageclass", "s", "", "StorageClass for PVC")
	pvcCmd.Flags().StringVarP(&o.Size, "size", "", "10", "Size of PVC in GBs")
	pvcCmd.Flags().StringVarP(&o.Namespace, "namespace", "", "default", "Namespace of PVC")

	o.PrintFlags.AddFlags(pvcCmd)
	cmdutil.AddDryRunFlag(pvcCmd)

	return pvcCmd
}

func (o *CreatePVCOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	name, err := cmdcreate.NameFromCommandArgs(cmd, args)
	if err != nil {
		return err
	}
	o.Name = name

	clientConfig, err := f.ToRESTConfig()
	if err != nil {
		return err
	}
	o.Client, err = corev1client.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	o.DryRunStrategy, err = cmdutil.GetDryRunStrategy(cmd)
	if err != nil {
		return err
	}
	cmdutil.PrintFlagsWithDryRunStrategy(o.PrintFlags, o.DryRunStrategy)

	printer, err := o.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	o.PrintObj = func(obj runtime.Object) error {
		return printer.PrintObj(obj, o.Out)
	}

	return nil
}

func (o *CreatePVCOptions) createPVCObject() *corev1.PersistentVolumeClaim {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.Name,
			Namespace: o.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					"storage": resource.MustParse(fmt.Sprintf("%sGi", o.Size)),
				},
			},
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
		},
	}

	if o.StorageClass != "" {
		pvc.Spec.StorageClassName = &o.StorageClass
	}

	return pvc
}

func (o *CreatePVCOptions) Run() error {
	pvc := o.createPVCObject()
	if o.DryRunStrategy != cmdutil.DryRunClient {
		createOptions := metav1.CreateOptions{}

		if o.DryRunStrategy == cmdutil.DryRunServer {
			createOptions.DryRun = []string{metav1.DryRunAll}
		}
		var err error
		pvc, err = o.Client.PersistentVolumeClaims(o.Namespace).Create(context.TODO(), pvc, createOptions)
		if err != nil {
			return fmt.Errorf("failed to create PVC: %v", err)
		}
	}

	return o.PrintObj(pvc)
}
