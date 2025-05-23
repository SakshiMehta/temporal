package calculator

var (
	_ Calculator          = (*ClusterAwareQuotaCalculator)(nil)
	_ NamespaceCalculator = (*ClusterAwareNamespaceQuotaCalculator)(nil)
)

type (
	// MemberCounter returns the total number of instances there are for a given service.
	MemberCounter interface {
		AvailableMemberCount() int
	}
	// ClusterAwareQuotaCalculator calculates the available quota for the current host based on the per instance and per
	// cluster quota. The quota could represent requests per second, total number of active requests, etc. It works by
	// dividing the per cluster quota by the total number of instances running the same service.
	ClusterAwareQuotaCalculator quotaCalculator[func() int]
	// ClusterAwareNamespaceQuotaCalculator is similar to ClusterAwareQuotaCalculator, but it uses quotas that
	// are specific to a namespace.
	ClusterAwareNamespaceQuotaCalculator quotaCalculator[func(namespace string) int]
	// quotaCalculator is a generic type that we use because the quota functions could be namespace specific or not.
	quotaCalculator[T any] struct {
		MemberCounter MemberCounter
		// PerInstanceQuota is a function that returns the per instance limit.
		PerInstanceQuota T
		// GlobalQuota is a function that returns the per cluster limit.
		GlobalQuota T
	}
)

// getQuota returns the effective resource limit for a host given the per instance and per cluster
// limits. The "resource" here could be requests per second, total number of active requests, etc. The cluster-wide
// limit is used if and only if it is configured to a value greater than zero and the number of instances that
// the memberCounter reports is greater than zero. Otherwise, the per-instance limit is used.
func getQuota(memberCounter MemberCounter, instanceLimit, clusterLimit int) float64 {
	if clusterLimit > 0 && memberCounter != nil {
		if clusterSize := memberCounter.AvailableMemberCount(); clusterSize > 0 {
			return float64(clusterLimit) / float64(clusterSize)
		}
	}

	return float64(instanceLimit)
}

func (l ClusterAwareQuotaCalculator) GetQuota() float64 {
	return getQuota(l.MemberCounter, l.PerInstanceQuota(), l.GlobalQuota())
}

func (l ClusterAwareNamespaceQuotaCalculator) GetQuota(namespace string) float64 {
	return getQuota(l.MemberCounter, l.PerInstanceQuota(namespace), l.GlobalQuota(namespace))
}
