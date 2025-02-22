// Copyright 2020-2021 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"os"
	"time"

	"github.com/cilium/cilium-cli/defaults"
	"github.com/cilium/cilium-cli/hubble"

	"github.com/spf13/cobra"
)

func newCmdHubble() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hubble",
		Short: "Hubble observability",
		Long:  ``,
	}

	cmd.AddCommand(
		newCmdHubbleEnable(),
		newCmdHubbleDisable(),
		newCmdPortForwardCommand(),
		newCmdUI(),
	)

	return cmd
}

func newCmdHubbleEnable() *cobra.Command {
	var params = hubble.Parameters{
		Writer: os.Stdout,
	}

	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable Hubble observability",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			h := hubble.NewK8sHubble(k8sClient, params)
			if err := h.Enable(context.Background()); err != nil {
				fatalf("Unable to enable Hubble:  %s", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&params.Namespace, "namespace", "n", "kube-system", "Namespace Cilium is running in")
	cmd.Flags().BoolVar(&params.Relay, "relay", true, "Deploy Hubble Relay")
	cmd.Flags().StringVar(&params.RelayImage, "relay-image", defaults.RelayImage, "Image path to use for Relay")
	cmd.Flags().StringVar(&params.RelayVersion, "relay-version", defaults.Version, "Version of Relay to deploy")
	cmd.Flags().StringVar(&params.RelayServiceType, "relay-service-type", "ClusterIP", "Type of Kubernetes service to expose Hubble Relay")
	cmd.Flags().BoolVar(&params.UI, "ui", false, "Enable Hubble UI")
	cmd.Flags().BoolVar(&params.CreateCA, "create-ca", false, "Automatically create CA if needed")
	cmd.Flags().StringVar(&contextName, "context", "", "Kubernetes configuration context")
	cmd.Flags().DurationVar(&params.CiliumReadyTimeout, "cilium-ready-timeout", 5*time.Minute,
		"Timeout for Cilium to become ready before deploying Hubble components")

	return cmd
}

func newCmdHubbleDisable() *cobra.Command {
	var params = hubble.Parameters{
		Writer: os.Stdout,
	}

	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable Hubble observability",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			h := hubble.NewK8sHubble(k8sClient, params)
			if err := h.Disable(context.Background()); err != nil {
				fatalf("Unable to disable Hubble:  %s", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&params.Namespace, "namespace", "n", "kube-system", "Namespace Cilium is running in")
	cmd.Flags().StringVar(&contextName, "context", "", "Kubernetes configuration context")

	return cmd
}

func newCmdPortForwardCommand() *cobra.Command {
	var params = hubble.Parameters{
		Writer: os.Stdout,
	}

	cmd := &cobra.Command{
		Use:   "port-forward",
		Short: "Forward the relay port to the local machine",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := params.PortForwardCommand(context.Background()); err != nil {
				fatalf("Unable to port forward: %s", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&params.Namespace, "namespace", "n", "kube-system", "Namespace Cilium is running in")
	cmd.Flags().StringVar(&params.Context, "context", "", "Kubernetes configuration context")
	cmd.Flags().IntVar(&params.PortForward, "port-forward", 4245, "Local port to forward to")

	return cmd
}

func newCmdUI() *cobra.Command {
	var params = hubble.Parameters{
		Writer: os.Stdout,
	}

	cmd := &cobra.Command{
		Use:   "ui",
		Short: "Open the Hubble UI",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := params.UIPortForwardCommand(context.Background()); err != nil {
				fatalf("Unable to port forward: %s", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&params.Namespace, "namespace", "n", "kube-system", "Namespace Cilium is running in")
	cmd.Flags().StringVar(&params.Context, "context", "", "Kubernetes configuration context")
	cmd.Flags().IntVar(&params.UIPortForward, "port-forward", 12000, "Local port to use for the port forward")

	return cmd
}
