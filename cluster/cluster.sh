export KUBEVIRT_PROVIDER=${KUBEVIRT_PROVIDER:-'kind-1.23-sriov'}
export KUBEVIRTCI_TAG=2301061052-bf2db41

KUBEVIRTCI_REPO='https://github.com/kubevirt/kubevirtci.git'
# The CLUSTER_PATH var is used in cluster folder and points to the _kubevirtci where the cluster is deployed from.
CLUSTER_PATH=${CLUSTER_PATH:-"${PWD}/_kubevirtci/"}

function cluster::_get_repo() {
    git --git-dir "${CLUSTER_PATH}"/.git remote get-url origin
}

function cluster::_get_tag() {
    git -C "${CLUSTER_PATH}" describe --tags
}

function cluster::install() {
    # Remove cloned kubevirtci repository if it does not match the requested one
    if [ -d "${CLUSTER_PATH}" ]; then
        if [ $(cluster::_get_repo) != ${KUBEVIRTCI_REPO} -o $(cluster::_get_tag) != ${KUBEVIRTCI_TAG} ]; then
          rm -rf "${CLUSTER_PATH}"
        fi
    fi

    if [ ! -d "${CLUSTER_PATH}" ]; then
        git clone https://github.com/kubevirt/kubevirtci.git "${CLUSTER_PATH}"
        (
            cd "${CLUSTER_PATH}"
            git checkout ${KUBEVIRTCI_TAG}
        )
    fi
}

function cluster::path() {
    echo -n "${CLUSTER_PATH}"
}

function cluster::kubeconfig() {
    echo -n "${CLUSTER_PATH}"/_ci-configs/"${KUBEVIRT_PROVIDER}"/.kubeconfig
}
