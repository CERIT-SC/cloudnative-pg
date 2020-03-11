/*
This file is part of Cloud Native PostgreSQL.

Copyright (C) 2019-2020 2ndQuadrant Italia SRL. Exclusively licensed to 2ndQuadrant Limited.
*/

package specs

import (
	"errors"
	"fmt"
	"strconv"

	corev1 "k8s.io/api/core/v1"
)

var (
	// ErrorContainerNotFound is raised when you are looking for the PostgreSQL
	// image in a Pod created by this operator but you don't find it
	ErrorContainerNotFound = errors.New("container not found")
)

// GetNodeSerial get the serial number of a Pod created by the operator
// for a Cluster
func GetNodeSerial(pod corev1.Pod) (int, error) {
	nodeSerial, ok := pod.Annotations[ClusterSerialAnnotationName]
	if !ok {
		return 0, fmt.Errorf("missing node serial annotation")
	}

	result, err := strconv.Atoi(nodeSerial)
	if err != nil {
		return 0, fmt.Errorf("wrong node serial annotation: %v", nodeSerial)
	}

	return result, nil
}

// IsPodPrimary check if a certain pod belongs to a primary
func IsPodPrimary(pod corev1.Pod) bool {
	role, hasRole := pod.ObjectMeta.Labels[ClusterRoleLabelName]
	if !hasRole {
		return false
	}

	return role == ClusterRoleLabelPrimary
}

// IsPodStandby check if a certain pod belongs to a standby
func IsPodStandby(pod corev1.Pod) bool {
	return !IsPodPrimary(pod)
}

// GetPostgreSQLImageName get the PostgreSQL image name used for this Pod
func GetPostgreSQLImageName(pod corev1.Pod) (string, error) {
	for _, container := range pod.Spec.Containers {
		if container.Name == PostgresContainerName {
			return container.Image, nil
		}
	}

	return "", ErrorContainerNotFound
}
