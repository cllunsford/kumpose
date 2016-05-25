# kumpose

Kumpose converts docker-compose.yml files into kubernetes configurations.

# Usage

```
# Convert 'docker-compose.yml' and print kubernetes JSON to stdout
$ kumpose

# Provide a path to the 'docker-compose.yml' and pipe to kubectl
$ kumpose -f path/to/docker-compose.yml | kubectl create -f -
```

# Mapping Logic

When there is no ambiguity, Kumpose maps docker-compose parameters directly to the equivalent Kubernetes parameter.

Kumpose defaults to mapping each docker-compose service to a Kubernetes `Deployment` object.  This will be configurable in the future to allow conversion to RC, Pod, Job, or DaemonSet.

For Kubernetes features that are not available in docker-compose, kumpose implements the following default values.  Future versions may allow these to be provided as flags when calling kumpose:

 * ImagePullPolicy - IfNotPresent
 * Replicas - 1

docker-compose parameters Not Yet Implemented:

 * aliases
 * cap\_add, cap\_drop
 * cgroup\_parent
 * container\_name
 * cpu\_quota
 * cpu\_shares
 * cpuset
 * devices
 * depends\_on
 * dns
 * dns\_search
 * domainname
 * env\_file
 * environment
 * expose
 * external\_links
 * extra\_hosts
 * hostname
 * ipc
 * labels
 * links
 * logging (v2)
 * log\_driver (v1)
 * log\_opt (v1)
 * mac\_address
 * mem\_limit
 * memswap\_limit
 * net (v1)
 * network\_mode (v2)
 * networks (v2)
 * pid
 * privileged
 * read\_only
 * restart
 * security\_opt
 * stop\_signal
 * stdin\_open
 * ulimits
 * user
 * volumes
 * volume\_driver
 * volumes\_from

# Notes

Currently importing k8s.io/kubernetes/pkg/api which is not designed for external consumption.  This increases the number of files imported and build times.  Will switch to versioned golang API once available (https://github.com/kubernetes/kubernetes/issues/5660).

# License

Chris Lunsford (c) 2016 - MIT License
