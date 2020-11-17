package clusterpolicy

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	gpuv1 "github.com/NVIDIA/gpu-operator/pkg/apis/nvidia/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type controlFunc []func(n ClusterPolicyController) (gpuv1.State, error)

func ServiceAccount(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].ServiceAccount.DeepCopy()
	logger := log.WithValues("ServiceAccount", obj.Name, "Namespace", obj.Namespace)

	if n.singleton.Status.StateRollback  != -1 {
		return gpuv1.Ready, nil
	}

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			//logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}

func Role(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].Role.DeepCopy()
	logger := log.WithValues("Role", obj.Name, "Namespace", obj.Namespace)

	if n.singleton.Status.StateRollback  != -1 {
		return gpuv1.Ready, nil
	}

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			//logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}

func RoleBinding(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].RoleBinding.DeepCopy()
	logger := log.WithValues("RoleBinding", obj.Name, "Namespace", obj.Namespace)

	if n.singleton.Status.StateRollback  != -1 {
		return gpuv1.Ready, nil
	}

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			//logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}

func ClusterRole(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].ClusterRole.DeepCopy()
	logger := log.WithValues("ClusterRole", obj.Name, "Namespace", obj.Namespace)

	if n.singleton.Status.StateRollback  != -1 {
		return gpuv1.Ready, nil
	}

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			//logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}

func ClusterRoleBinding(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].ClusterRoleBinding.DeepCopy()
	logger := log.WithValues("ClusterRoleBinding", obj.Name, "Namespace", obj.Namespace)

	if n.singleton.Status.StateRollback  != -1 {
		return gpuv1.Ready, nil
	}

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			//logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}

func ConfigMap(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].ConfigMap.DeepCopy()
	logger := log.WithValues("ConfigMap", obj.Name, "Namespace", obj.Namespace)

	if n.singleton.Status.StateRollback  != -1 {
		return gpuv1.Ready, nil
	}

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	// delete to simplify dev-time modifications
	n.rec.client.Delete(context.TODO(), obj)

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			//logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}

func kernelFullVersion(n ClusterPolicyController) (string, string, string) {
	logger := log.WithValues("Request.Namespace", "default", "Request.Name", "Node")
	// We need the node labels to fetch the correct container
	opts := []client.ListOption{
		client.MatchingLabels{"nvidia.com/gpu.present": "true"},
	}

	list := &corev1.NodeList{}
	err := n.rec.client.List(context.TODO(), list, opts...)
	if err != nil {
		logger.Info("Could not get NodeList", "ERROR", err)
		return "", "", ""
	}

	if len(list.Items) == 0 {
		// none of the nodes matched nvidia GPU label
		// either the nodes do not have GPUs, or NFD is not running
		logger.Info("Could not get any nodes to match nvidia.com/gpu.present label", "ERROR", "")
		return "", "", ""
	}

	// Assuming all nodes are running the same kernel version,
	// One could easily add driver-kernel-versions for each node.
	node := list.Items[0]
	labels := node.GetLabels()

	var ok bool
	kFVersion, ok := labels["feature.node.kubernetes.io/kernel-version.full"]
	if ok {
		//logger.Info(kFVersion)
	} else {
		err := errors.NewNotFound(schema.GroupResource{Group: "Node", Resource: "Label"},
			"feature.node.kubernetes.io/kernel-version.full")
		logger.Info("Couldn't get kernelVersion, did you run the node feature discovery?", err)
		return "", "", ""
	}

	osName, ok := labels["feature.node.kubernetes.io/system-os_release.ID"]
	if !ok {
		return kFVersion, "", ""
	}
	osVersion, ok := labels["feature.node.kubernetes.io/system-os_release.VERSION_ID"]
	if !ok {
		return kFVersion, "", ""
	}
	osTag := fmt.Sprintf("%s%s", osName, osVersion)

	return kFVersion, osTag, osVersion
}

func getDcgmExporter() string {
	dcgmExporter := os.Getenv("NVIDIA_DCGM_EXPORTER")
	if dcgmExporter == "" {
		log.Info(fmt.Sprintf("ERROR: Could not find environment variable NVIDIA_DCGM_EXPORTER"))
		os.Exit(1)
	}
	return dcgmExporter
}

func preProcessDaemonSet(obj *appsv1.DaemonSet, n ClusterPolicyController) {
	transformations := map[string]func(*appsv1.DaemonSet, *gpuv1.ClusterPolicySpec, ClusterPolicyController) error{
		"nvidia-driver-daemonset":            TransformDriver,
		"nvidia-mig-mode-daemonset":          TransformMigMode,
		"nvidia-container-toolkit-daemonset": TransformToolkit,
		"nvidia-device-plugin-daemonset":     TransformDevicePlugin,
		"nvidia-dcgm-exporter":               TransformDCGMExporter,
		"gpu-feature-discovery":              TransformGPUDiscoveryPlugin,
	}

	t, ok := transformations[obj.Name]
	if !ok {
		log.Info(fmt.Sprintf("No transformation for Daemonset '%s'", obj.Name))
		return
	}

	err := t(obj, &n.singleton.Spec, n)
	if err != nil {
		log.Info(fmt.Sprintf("FATAL: Failed to apply transformation '%s' with error: '%v'", obj.Name, err))
		os.Exit(1)
	}
}

// TransformGPUDiscoveryPlugin transforms GPU discovery daemonset with required config as per ClusterPolicy
func TransformGPUDiscoveryPlugin(obj *appsv1.DaemonSet, config *gpuv1.ClusterPolicySpec, n ClusterPolicyController) error {
	// update image
	obj.Spec.Template.Spec.Containers[0].Image = config.GroupFeatureDiscovery.ImagePath()

	// update image pull policy
	if config.GroupFeatureDiscovery.ImagePullPolicy != "" {
		obj.Spec.Template.Spec.Containers[0].ImagePullPolicy = config.GroupFeatureDiscovery.ImagePolicy(config.GroupFeatureDiscovery.ImagePullPolicy)
	}

	// set image pull secrets
	if len(config.GroupFeatureDiscovery.ImagePullSecrets) > 0 {
		for _, secret := range config.GroupFeatureDiscovery.ImagePullSecrets {
			obj.Spec.Template.Spec.ImagePullSecrets = append(obj.Spec.Template.Spec.ImagePullSecrets, v1.LocalObjectReference{Name: secret})
		}
	}

	// set node selector if specified
	if len(config.GroupFeatureDiscovery.NodeSelector) > 0 {
		obj.Spec.Template.Spec.NodeSelector = config.GroupFeatureDiscovery.NodeSelector
	}

	// set node affinity if specified
	if config.GroupFeatureDiscovery.Affinity != nil {
		obj.Spec.Template.Spec.Affinity = config.GroupFeatureDiscovery.Affinity
	}

	// set tolerations if specified
	if len(config.GroupFeatureDiscovery.Tolerations) > 0 {
		obj.Spec.Template.Spec.Tolerations = config.GroupFeatureDiscovery.Tolerations
	}

	// set resource limits
	if config.GroupFeatureDiscovery.Resources != nil {
		// apply resource limits to all containers
		for i := range obj.Spec.Template.Spec.Containers {
			obj.Spec.Template.Spec.Containers[i].Resources = *config.GroupFeatureDiscovery.Resources
		}
	}

	// update MIG strategy and discovery intervals
	var migStrategy gpuv1.MigStrategy = gpuv1.MigStrategyNone
	if config.GroupFeatureDiscovery.MigStrategy != "" {
		migStrategy = config.GroupFeatureDiscovery.MigStrategy
	}
	setContainerEnv(&(obj.Spec.Template.Spec.Containers[0]), "GFD_MIG_STRATEGY", fmt.Sprintf("%s", migStrategy))
	if migStrategy != gpuv1.MigStrategyNone {
		setContainerEnv(&(obj.Spec.Template.Spec.Containers[0]), "NVIDIA_MIG_MONITOR_DEVICES", "all")
	}

	// update discovery interval
	discoveryIntervalSeconds := 60
	if config.GroupFeatureDiscovery.DiscoveryIntervalSeconds != 0 {
		discoveryIntervalSeconds = config.GroupFeatureDiscovery.DiscoveryIntervalSeconds
	}
	setContainerEnv(&(obj.Spec.Template.Spec.Containers[0]), "GFD_SLEEP_INTERVAL", fmt.Sprintf("%ds", discoveryIntervalSeconds))

	return nil
}

// Read and parse os-release file
func parseOSRelease() (map[string]string, error) {
	release := map[string]string{}

	f, err := os.Open("/host-etc/os-release")
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`^(?P<key>\w+)=(?P<value>.+)`)

	// Read line-by-line
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if m := re.FindStringSubmatch(line); m != nil {
			release[m[1]] = strings.Trim(m[2], `"`)
		}
	}

	return release, nil
}

// TransformDriver transforms Nvidia driver daemonset with required config as per ClusterPolicy
func TransformDriver(obj *appsv1.DaemonSet, config *gpuv1.ClusterPolicySpec, n ClusterPolicyController) error {
	kvers, osTag, _ := kernelFullVersion(n)
	if kvers == "" {
		return fmt.Errorf("ERROR: Could not find kernel full version: ('%s', '%s')", kvers, osTag)
	}

	img := fmt.Sprintf("%s-%s", config.Driver.ImagePath(), osTag)
	obj.Spec.Template.Spec.Containers[0].Image = img

	// update image pull policy
	if config.Driver.ImagePullPolicy != "" {
		obj.Spec.Template.Spec.Containers[0].ImagePullPolicy = config.Driver.ImagePolicy(config.Driver.ImagePullPolicy)
	}
	// set image pull secrets
	if len(config.Driver.ImagePullSecrets) > 0 {
		for _, secret := range config.Driver.ImagePullSecrets {
			obj.Spec.Template.Spec.ImagePullSecrets = append(obj.Spec.Template.Spec.ImagePullSecrets, v1.LocalObjectReference{Name: secret})
		}
	}
	// set node selector if specified
	if len(config.Driver.NodeSelector) > 0 {
		obj.Spec.Template.Spec.NodeSelector = config.Driver.NodeSelector
	}
	// set node affinity if specified
	if config.Driver.Affinity != nil {
		obj.Spec.Template.Spec.Affinity = config.Driver.Affinity
	}
	// set tolerations if specified
	if len(config.Driver.Tolerations) > 0 {
		obj.Spec.Template.Spec.Tolerations = config.Driver.Tolerations
	}
	// set resource limits
	if config.Driver.Resources != nil {
		// apply resource limits to all containers
		for i := range obj.Spec.Template.Spec.Containers {
			obj.Spec.Template.Spec.Containers[i].Resources = *config.Driver.Resources
		}
	}

	// Inject EUS kernel RPM's as an override to the entrypoint
	// Add Env Vars needed by nvidia-driver to enable the right releasever and rpm repo
	if !strings.Contains(osTag, "rhel") && !strings.Contains(osTag, "rhcos") {
		return nil
	}

	release, err := parseOSRelease()
	if err != nil {
		return fmt.Errorf("ERROR: failed to get os-release: %s", err)
	}

	ocpV, err := OpenshiftVersion()
	if err != nil {
		// might be RHEL node using upstream k8s, don't error out.
		log.Info(fmt.Sprintf("ERROR: failed to get OpenShift version: %s", err))
	}

	rhelVersion := corev1.EnvVar{Name: "RHEL_VERSION", Value: release["RHEL_VERSION"]}
	ocpVersion := corev1.EnvVar{Name: "OPENSHIFT_VERSION", Value: ocpV}

	obj.Spec.Template.Spec.Containers[0].Env = append(obj.Spec.Template.Spec.Containers[0].Env, rhelVersion)
	obj.Spec.Template.Spec.Containers[0].Env = append(obj.Spec.Template.Spec.Containers[0].Env, ocpVersion)

	return nil
}

// TransformMigMode transforms Nvidia mig-mode daemonset with required config as per ClusterPolicy
func TransformMigMode(obj *appsv1.DaemonSet, config *gpuv1.ClusterPolicySpec, n ClusterPolicyController) error {
	if err := TransformDriver(obj, config, n); err != nil {
		return err
	}

	migStrategy := corev1.EnvVar{
		Name: "DRV_MIG_STRATEGY",
		Value: string(n.singleton.Spec.GroupFeatureDiscovery.MigStrategy),
	}
	migMode := corev1.EnvVar{
		Name: "DRV_MIG_MODE",
		Value: n.singleton.Spec.Driver.MigMode,
	}

	obj.Spec.Template.Spec.Containers[0].Env = append(obj.Spec.Template.Spec.Containers[0].Env, migStrategy)
	obj.Spec.Template.Spec.Containers[0].Env = append(obj.Spec.Template.Spec.Containers[0].Env, migMode)

	return nil
}

// TransformToolkit transforms Nvidia container-toolkit daemonset with required config as per ClusterPolicy
func TransformToolkit(obj *appsv1.DaemonSet, config *gpuv1.ClusterPolicySpec, n ClusterPolicyController) error {
	obj.Spec.Template.Spec.Containers[0].Image = config.Toolkit.ImagePath()
	runtime := string(config.Operator.DefaultRuntime)

	setContainerEnv(&(obj.Spec.Template.Spec.Containers[0]), "RUNTIME", runtime)
	if runtime == "docker" {
		setContainerEnv(&(obj.Spec.Template.Spec.Containers[0]), "RUNTIME_ARGS",
			"--socket /var/run/docker.sock")
	}

	// update image pull policy
	if config.Toolkit.ImagePullPolicy != "" {
		obj.Spec.Template.Spec.Containers[0].ImagePullPolicy = config.Toolkit.ImagePolicy(config.Toolkit.ImagePullPolicy)
	}
	// set image pull secrets
	if len(config.Toolkit.ImagePullSecrets) > 0 {
		for _, secret := range config.Toolkit.ImagePullSecrets {
			obj.Spec.Template.Spec.ImagePullSecrets = append(obj.Spec.Template.Spec.ImagePullSecrets, v1.LocalObjectReference{Name: secret})
		}
	}
	// set node selector if specified
	if len(config.Toolkit.NodeSelector) > 0 {
		obj.Spec.Template.Spec.NodeSelector = config.Toolkit.NodeSelector
	}
	// set node affinity if specified
	if config.Toolkit.Affinity != nil {
		obj.Spec.Template.Spec.Affinity = config.Toolkit.Affinity
	}
	// set tolerations if specified
	if len(config.Toolkit.Tolerations) > 0 {
		obj.Spec.Template.Spec.Tolerations = config.Toolkit.Tolerations
	}
	// set resource limits
	if config.Toolkit.Resources != nil {
		// apply resource limits to all containers
		for i := range obj.Spec.Template.Spec.Containers {
			obj.Spec.Template.Spec.Containers[i].Resources = *config.Toolkit.Resources
		}
	}

	return nil
}

// TransformDevicePlugin transforms k8s-device-plugin daemonset with required config as per ClusterPolicy
func TransformDevicePlugin(obj *appsv1.DaemonSet, config *gpuv1.ClusterPolicySpec, n ClusterPolicyController) error {
	// update image
	obj.Spec.Template.Spec.Containers[0].Image = config.DevicePlugin.ImagePath()
	// update image pull policy
	if config.DevicePlugin.ImagePullPolicy != "" {
		obj.Spec.Template.Spec.Containers[0].ImagePullPolicy = config.DevicePlugin.ImagePolicy(config.DevicePlugin.ImagePullPolicy)
	}
	// set image pull secrets
	if len(config.DevicePlugin.ImagePullSecrets) > 0 {
		for _, secret := range config.DevicePlugin.ImagePullSecrets {
			obj.Spec.Template.Spec.ImagePullSecrets = append(obj.Spec.Template.Spec.ImagePullSecrets, v1.LocalObjectReference{Name: secret})
		}
	}
	// set node selector if specified
	if len(config.DevicePlugin.NodeSelector) > 0 {
		obj.Spec.Template.Spec.NodeSelector = config.DevicePlugin.NodeSelector
	}
	// set node affinity if specified
	if config.DevicePlugin.Affinity != nil {
		obj.Spec.Template.Spec.Affinity = config.DevicePlugin.Affinity
	}
	// set tolerations if specified
	if len(config.DevicePlugin.Tolerations) > 0 {
		obj.Spec.Template.Spec.Tolerations = config.DevicePlugin.Tolerations
	}
	// set resource limits
	if config.DevicePlugin.Resources != nil {
		// apply resource limits to all containers
		for i := range obj.Spec.Template.Spec.Containers {
			obj.Spec.Template.Spec.Containers[i].Resources = *config.DevicePlugin.Resources
		}
	}

	var migStrategy gpuv1.MigStrategy = config.GroupFeatureDiscovery.MigStrategy
	if migStrategy != "" {
		addContainerArg(&(obj.Spec.Template.Spec.Containers[0]), "-mig-strategy")
		addContainerArg(&(obj.Spec.Template.Spec.Containers[0]), fmt.Sprintf("%s", migStrategy))
	}
	return nil
}

// TransformDCGMExporter transforms dcgm exporter daemonset with required config as per ClusterPolicy
func TransformDCGMExporter(obj *appsv1.DaemonSet, config *gpuv1.ClusterPolicySpec, n ClusterPolicyController) error {
	// update image
	obj.Spec.Template.Spec.Containers[0].Image = config.DCGMExporter.ImagePath()

	// update image pull policy
	if config.DCGMExporter.ImagePullPolicy != "" {
		obj.Spec.Template.Spec.Containers[0].ImagePullPolicy = config.DevicePlugin.ImagePolicy(config.DCGMExporter.ImagePullPolicy)
	}
	// set image pull secrets
	if len(config.DCGMExporter.ImagePullSecrets) > 0 {
		for _, secret := range config.DCGMExporter.ImagePullSecrets {
			obj.Spec.Template.Spec.ImagePullSecrets = append(obj.Spec.Template.Spec.ImagePullSecrets, v1.LocalObjectReference{Name: secret})
		}
	}
	// set node selector if specified
	if len(config.DCGMExporter.NodeSelector) > 0 {
		obj.Spec.Template.Spec.NodeSelector = config.DCGMExporter.NodeSelector
	}
	// set node affinity if specified
	if config.DCGMExporter.Affinity != nil {
		obj.Spec.Template.Spec.Affinity = config.DCGMExporter.Affinity
	}
	// set tolerations if specified
	if len(config.DCGMExporter.Tolerations) > 0 {
		obj.Spec.Template.Spec.Tolerations = config.DCGMExporter.Tolerations
	}
	// set resource limits
	if config.DCGMExporter.Resources != nil {
		// apply resource limits to all containers
		for i := range obj.Spec.Template.Spec.Containers {
			obj.Spec.Template.Spec.Containers[i].Resources = *config.DCGMExporter.Resources
		}
	}

	kvers, osTag, _ := kernelFullVersion(n)
	if kvers == "" {
		return fmt.Errorf("ERROR: Could not find kernel full version: ('%s', '%s')", kvers, osTag)
	}

	if !strings.Contains(osTag, "rhel") && !strings.Contains(osTag, "rhcos") {
		return nil
	}

	// update init container config for per pod specific resources
	initContainerImage, initContainerName, initContainerCmd := "ubi8/ubi-minimal:8.2-349", "init-pod-nvidia-metrics-exporter", "/bin/entrypoint.sh"
	initContainer := v1.Container{}
	obj.Spec.Template.Spec.InitContainers = append(obj.Spec.Template.Spec.InitContainers, initContainer)
	obj.Spec.Template.Spec.InitContainers[0].Image = initContainerImage
	obj.Spec.Template.Spec.InitContainers[0].Name = initContainerName
	obj.Spec.Template.Spec.InitContainers[0].Command = []string{initContainerCmd}

	// need CAP_SYS_ADMIN privileges for collecting pod specific resources
	privileged := true
	securityContext := &corev1.SecurityContext{
		Privileged: &privileged,
	}

	// Add initContainer for OCP to set proper SELinux context on /var/lib/kubelet/pod-resources
	obj.Spec.Template.Spec.InitContainers[0].SecurityContext = securityContext

	volMountSockName, volMountSockPath := "pod-gpu-resources", "/var/lib/kubelet/pod-resources"
	volMountSock := corev1.VolumeMount{Name: volMountSockName, MountPath: volMountSockPath}
	obj.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(obj.Spec.Template.Spec.InitContainers[0].VolumeMounts, volMountSock)

	volMountConfigName, volMountConfigPath, volMountConfigSubPath := "init-config", "/bin/entrypoint.sh", "entrypoint.sh"
	volMountConfig := corev1.VolumeMount{Name: volMountConfigName, ReadOnly: true, MountPath: volMountConfigPath, SubPath: volMountConfigSubPath}
	obj.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(obj.Spec.Template.Spec.InitContainers[0].VolumeMounts, volMountConfig)

	volMountConfigKey, volMountConfigDefaultMode := "nvidia-dcgm-exporter", int32(0700)
	initVol := corev1.Volume{Name: volMountConfigName, VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: volMountConfigKey}, DefaultMode: &volMountConfigDefaultMode}}}
	obj.Spec.Template.Spec.Volumes = append(obj.Spec.Template.Spec.Volumes, initVol)

	return nil
}

func setContainerEnv(c *corev1.Container, key, value string) {
	for i, val := range c.Env {
		if val.Name != key {
			continue
		}

		c.Env[i].Value = value
		return
	}

	//log.Info(fmt.Sprintf("Info: Could not find environment variable %s in container %s, appending it", key, c.Name))
	c.Env = append(c.Env, corev1.EnvVar{Name: key, Value: value})
}

func addContainerArg(c *corev1.Container, val string) {
	c.Args = append(c.Args, val)
	log.Info(fmt.Sprintf("Info: adding Pod argument '%s'", val))
}

func TransformDevicePluginValidator(obj *v1.Pod, config *gpuv1.ClusterPolicySpec, n ClusterPolicyController) (gpuv1.State, error) {
	if state, err := TransformValidator(obj, config, n); err != nil {
		return state, err
	}

	if config.GroupFeatureDiscovery.MigStrategy != "mixed" {
		return gpuv1.Ready, nil
	}

	opts := []client.ListOption{
		client.MatchingLabels{"nvidia.com/gpu.present": "true"},
	}
	log.Info("DEBUG: Node", "NodeSelector", "nvidia.com/gpu.present=true")
	list := &corev1.NodeList{}
	err := n.rec.client.List(context.TODO(), list, opts...)

	if err != nil {
		log.Info("Could not get NodeList", err)
	}
	log.Info("DEBUG: Node", "NumberOfNodes", len(list.Items))
	if len(list.Items) == 0 {
		return gpuv1.Ready, fmt.Errorf("ERROR: Could not find any node with label 'nvidia.com/gpu.present=true'")
	}

	node := list.Items[0]
	log.Info("DEBUG: Node", "NodeName", node.ObjectMeta.Name)

	var mig_resource_name corev1.ResourceName = ""

	for resource_name, quantity := range node.Status.Capacity {
		if quantity.Value() < 1 || ! strings.HasPrefix(string(resource_name), "nvidia.com/mig-") {
			continue
		}

		mig_resource_name = resource_name
		break
	}

	if mig_resource_name == "" {
		return gpuv1.NotReady, fmt.Errorf("ERROR: Could not find any 'nvidia.com/mig-*' resource in node/%s",
			              node.ObjectMeta.Name)
	}

	mig_resource := corev1.ResourceList{
		mig_resource_name: resource.MustParse("1"),
	}

	obj.Spec.Containers[0].Resources.Limits = mig_resource
	obj.Spec.Containers[0].Resources.Requests = mig_resource

	return gpuv1.Ready, nil
}

// TransformValidator transforms driver and device plugin validator pods with required config as per ClusterPolicy
func TransformValidator(obj *v1.Pod, config *gpuv1.ClusterPolicySpec, n ClusterPolicyController) (gpuv1.State, error) {
	// update image
	if config.Operator.Validator.Repository != "" {
		obj.Spec.Containers[0].Image = config.Operator.Validator.ImagePath()
	}
	// update image pull policy
	if config.Operator.Validator.ImagePullPolicy != "" {
		obj.Spec.Containers[0].ImagePullPolicy = config.Operator.Validator.ImagePolicy(config.Operator.Validator.ImagePullPolicy)
	}
	// set image pull secrets
	if config.Operator.Validator.ImagePullSecrets != nil {
		for _, secret := range config.Operator.Validator.ImagePullSecrets {
			obj.Spec.ImagePullSecrets = append(obj.Spec.ImagePullSecrets, v1.LocalObjectReference{Name: secret})
		}
	}
	return gpuv1.Ready, nil
}

func isDeploymentReady(name string, n ClusterPolicyController) gpuv1.State {
	opts := []client.ListOption{
		client.MatchingLabels{"app": name},
	}
	//log.Info("DEBUG: DaemonSet", "LabelSelector", fmt.Sprintf("app=%s", name))
	list := &appsv1.DeploymentList{}
	err := n.rec.client.List(context.TODO(), list, opts...)
	if err != nil {
		log.Info("Could not get DaemonSetList", err)
	}
	//log.Info("DEBUG: DaemonSet", "NumberOfDaemonSets", len(list.Items))
	if len(list.Items) == 0 {
		return gpuv1.NotReady
	}

	ds := list.Items[0]
	//log.Info("DEBUG: DaemonSet", "NumberUnavailable", ds.Status.UnavailableReplicas)

	if ds.Status.UnavailableReplicas != 0 {
		return gpuv1.NotReady
	}

	return isPodReady(name, n, "Running")
}

func isDaemonSetReady(name string, n ClusterPolicyController) gpuv1.State {
	opts := []client.ListOption{
		client.MatchingLabels{"app": name},
	}
	//log.Info("DEBUG: DaemonSet", "LabelSelector", fmt.Sprintf("app=%s", name))
	list := &appsv1.DaemonSetList{}
	err := n.rec.client.List(context.TODO(), list, opts...)
	if err != nil {
		log.Info("Could not get DaemonSetList", err)
	}

	if len(list.Items) == 0 {
		log.Info("DEBUG: DaemonSet", "NumberOfDaemonSets", len(list.Items))
		return gpuv1.NotReady
	}

	ds := list.Items[0]

	if ds.Status.NumberUnavailable != 0 {
		log.Info("DEBUG: DaemonSet", "NumberUnavailable", ds.Status.NumberUnavailable)
		return gpuv1.NotReady
	}

	return isPodReady(name, n, "Running")
}

func Deployment(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].Deployment.DeepCopy()

	if n.singleton.Status.StateRollback  != -1 {
		return gpuv1.Ready, nil
	}

	logger := log.WithValues("Deployment", obj.Name, "Namespace", obj.Namespace)

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			//logger.Info("Found Resource")
			return isDeploymentReady(obj.Name, n), nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return isDeploymentReady(obj.Name, n), nil
}

func DaemonSet(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].DaemonSet.DeepCopy()

	if n.singleton.Status.StateRollback == state {

		return deleteDaemonSet(n, obj)
	} else if n.singleton.Status.StateRollback  != -1 {
		return gpuv1.Ready, nil
	}

	log.Info("STEP start", "name", obj.Name)
	if obj.Name == "nvidia-device-plugin-daemonset" {
		if n.singleton.Status.MigStrategy != "" && n.singleton.Status.MigStrategy != n.singleton.Spec.GroupFeatureDiscovery.MigStrategy {
			log.Info("ROLLBACK strategy-trigger", "old", n.singleton.Status.MigStrategy, "new", n.singleton.Spec.GroupFeatureDiscovery.MigStrategy)
			// the mig strategy changed, rollback the following states
			n.singleton.Status.StateRollback = state
			log.Info("ROLLBACK start", "StateRollback", ctrl.singleton.Status.StateRollback)

			n.singleton.Status.MigStrategy = ""

			return deleteDaemonSet(n, obj)
		} else {
			n.singleton.Status.MigStrategy = n.singleton.Spec.GroupFeatureDiscovery.MigStrategy
			if n.singleton.Status.MigStrategy == "" {
				n.singleton.Status.MigStrategy = "none"
			}
		}
	}

	if obj.Name == "nvidia-mig-mode-daemonset" {
		if n.singleton.Status.MigMode != "" && n.singleton.Status.MigMode != n.singleton.Spec.Driver.MigMode {
			log.Info("ROLLBACK mode-trigger", "old", n.singleton.Status.MigMode, "new", n.singleton.Spec.Driver.MigMode)
			// the mig mode changed, rollback the following states
			n.singleton.Status.StateRollback = state
			log.Info("ROLLBACK start", "StateRollback", ctrl.singleton.Status.StateRollback)
			n.singleton.Status.MigMode = ""

			return deleteDaemonSet(n, obj)
		} else {
			n.singleton.Status.MigMode = n.singleton.Spec.Driver.MigMode
			if n.singleton.Status.MigMode == "" {
				n.singleton.Status.MigMode = "none"
			}
		}
	}

	preProcessDaemonSet(obj, n)
	logger := log.WithValues("DaemonSet", obj.Name, "Namespace", obj.Namespace)

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			logger.Info("Resource Already Exists")
			return isDaemonSetReady(obj.Name, n), nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}
	logger.Info("Resource Created")

	return isDaemonSetReady(obj.Name, n), nil
}

func deleteDaemonSet(n ClusterPolicyController, obj *appsv1.DaemonSet) (gpuv1.State, error) {
	log.Info("ROLLBACK, Delete", "DaemonSet", obj.ObjectMeta.Name)
	err := n.rec.client.Delete(context.TODO(), obj);

	if errors.IsNotFound(err) {
		log.Info("Info: Delete DaemonSet", "err", err)

		// not an error, ignore
		err = nil
	} else {
		log.Info("Warning: Delete DaemonSet", "err", err)

	}


	return gpuv1.Ready, err
}

// The operator starts two pods in different stages to validate
// the correct working of the DaemonSets (driver and dp). Therefore
// the operator waits until the Pod completes and checks the error status
// to advance to the next state.
func isPodReady(name string, n ClusterPolicyController, phase corev1.PodPhase) gpuv1.State {
	opts := []client.ListOption{&client.MatchingLabels{"app": name}}

	//log.Info("DEBUG: Pod", "LabelSelector", fmt.Sprintf("app=%s", name))
	list := &corev1.PodList{}
	err := n.rec.client.List(context.TODO(), list, opts...)
	if err != nil {
		log.Info("Could not get PodList", err)
	}
	//log.Info("DEBUG: Pod", "NumberOfPods", len(list.Items))
	if len(list.Items) == 0 {
		log.Info("DEBUG: Pod not ready", "NumberOfPods", 0)
		return gpuv1.NotReady
	}

	pd := list.Items[0]

	if pd.Status.Phase != phase {
		log.Info("DEBUG: Pod not ready", "Phase", pd.Status.Phase, "!=", phase)
		return gpuv1.NotReady
	}
	//log.Info("DEBUG: Pod", "Phase", pd.Status.Phase, "==", phase)
	return gpuv1.Ready
}

func Pod(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].Pod.DeepCopy()
	log.Info("STEP start", "name", obj.Name)

	if n.singleton.Status.StateRollback == state {
		return deletePod(n, obj)
	} else if n.singleton.Status.StateRollback  != -1 {
		return gpuv1.Ready, nil
	}

	if state, err := preProcessPod(obj, n); err != nil {
		return state, err
	}
	logger := log.WithValues("Pod", obj.Name, "Namespace", obj.Namespace)

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			logger.Info("Resource Already Exists")
			return isPodReady(obj.Name, n, "Succeeded"), nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}
	logger.Info("Resource Created")

	return isPodReady(obj.Name, n, "Succeeded"), nil
}

func deletePod(n ClusterPolicyController, obj *corev1.Pod) (gpuv1.State, error) {
	log.Info("DEBUG: Delete", "Pod", obj.ObjectMeta.Name)

	err := n.rec.client.Delete(context.TODO(), obj);

	if errors.IsNotFound(err) {
		// not an error, ignore
		err = nil
	}

	return gpuv1.Ready, err
}

func preProcessPod(obj *v1.Pod, n ClusterPolicyController) (gpuv1.State, error) {
	transformations := map[string]func(*v1.Pod, *gpuv1.ClusterPolicySpec, ClusterPolicyController)  (gpuv1.State, error) {
		"nvidia-driver-validation":        TransformValidator,
		"nvidia-device-plugin-validation": TransformDevicePluginValidator,
	}

	t, ok := transformations[obj.Name]
	if !ok {
		log.Info(fmt.Sprintf("No transformation for Pod '%s'", obj.Name))
		return gpuv1.Ready, nil
	}

	state, err := t(obj, &n.singleton.Spec, n)
	if err != nil {
		if state == gpuv1.NotReady {
			log.Info(fmt.Sprintf("Failed to apply transformation '%s' with error: '%v'", obj.Name, err))
			return gpuv1.NotReady, nil
		}

		// permanent failure ...
		log.Info(fmt.Sprintf("FATAL: Failed to apply transformation '%s' with error: '%v'", obj.Name, err))
		os.Exit(1)
	}

	return gpuv1.Ready, nil
}

func SecurityContextConstraints(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].SecurityContextConstraints.DeepCopy()
	logger := log.WithValues("SecurityContextConstraints", obj.Name, "Namespace", "default")

	if n.singleton.Status.StateRollback  != -1 {
		return gpuv1.Ready, nil
	}

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			//logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}

func Service(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].Service.DeepCopy()
	logger := log.WithValues("Service", obj.Name, "Namespace", obj.Namespace)

	if n.singleton.Status.StateRollback  != -1 {
		return gpuv1.Ready, nil
	}

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			//logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}

func ServiceMonitor(n ClusterPolicyController) (gpuv1.State, error) {
	state := n.idx
	obj := n.resources[state].ServiceMonitor.DeepCopy()
	logger := log.WithValues("ServiceMonitor", obj.Name, "Namespace", obj.Namespace)

	if n.singleton.Status.StateRollback  != -1 {
		return gpuv1.Ready, nil
	}

	if err := controllerutil.SetControllerReference(n.singleton, obj, n.rec.scheme); err != nil {
		return gpuv1.NotReady, err
	}

	if err := n.rec.client.Create(context.TODO(), obj); err != nil {
		if errors.IsAlreadyExists(err) {
			//logger.Info("Found Resource")
			return gpuv1.Ready, nil
		}

		logger.Info("Couldn't create", "Error", err)
		return gpuv1.NotReady, err
	}

	return gpuv1.Ready, nil
}
