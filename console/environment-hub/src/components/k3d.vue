<script lang="ts" setup>
import { ref, reactive } from 'vue'
import { GetClusters, GetCluster, CreateClusters, DeleteCluster, StartCluster, StopCluster, InstallHelm } from './api'
import { Codemirror } from 'vue-codemirror'

interface Cluster {
  name: string
  servers: number
  agents: number
  port: number
  kubeconfig: string
  nodes: []
}

const tableData = ref([])

function loadClusters() {
  GetClusters((d) => {
    tableData.value = d
  })
}
loadClusters()

function openClusterCreationDialog() {
  creationDialogVisible.value = true
}

function createCluster() {
  creationDialogVisible.value = false
  CreateClusters(clusterForm, () => {
    loadClusters()
  })
}

function startCluster(name: string) {
  StartCluster(name, (d) => {
    loadClusters()
  })
}

function stopCluster(name: string) {
  StopCluster(name, (d) => {
    loadClusters()
  })
}

function deleteCluster(row: Cluster) {
  DeleteCluster(row.name, (d) => {
    loadClusters()
  })
}

const creationDialogVisible = ref(false)
const clusterForm = reactive({
  servers: 1,
  agents: 1,
  port: 30000
} as Cluster)
function openClusterEditor(cluster: Cluster) {
  clusterForm.servers = cluster.servers
  clusterForm.agents = cluster.agents
  clusterForm.port = cluster.port
  creationDialogVisible.value = true
}

const clusterDetail = ref(false)
const currentCluster = ref({} as Cluster)
function showClusterDetail(name: string) {
  clusterDetail.value = true
  GetCluster(name, (d) => {
    currentCluster.value = d
  })
}

const helmInstallDialog = ref(false)
function openHelmInstallDialog(name: string) {
  helmInstallDialog.value = true
  GetCluster(name, (d) => {
    currentCluster.value = d
  })
}
const helmForm = reactive({
  name: "",
  namespace: "default",
  repoURL: "",
  version: "",
  valuesStr: "",
  values: [""]
})
function installHelmChart() {
  const items = helmForm.valuesStr.split(",")
  for (var i=0;i<items.length;i++) {
    helmForm.values[i]=items[i]
  }
  InstallHelm(currentCluster.value.name, helmForm, (d) => {
  })
}
</script>

<template>
  <el-button @click="openClusterCreationDialog">Create</el-button>

  <el-table :data="tableData" stripe style="width: 100%">
    <el-table-column label="Name" width="180">
      <template #default="scope">
        <el-button link type="primary" size="small" @click="showClusterDetail(scope.row.name)">{{ scope.row.name }}</el-button>
      </template>
    </el-table-column>
    <el-table-column prop="nodes.length" label="Nodes"/>
    <el-table-column label="Port">
      <template #default="scope">
        <div v-for="(value, key) in scope.row.portBinding" :key="key">
          {{ key }}: {{ value }}
        </div>
      </template>
    </el-table-column>
    <el-table-column fixed="right" label="Operations" width="120">
      <template #default="scope">
        <el-button link type="primary" size="small" @click="openHelmInstallDialog(scope.row.name)">Install</el-button>
        <el-button link type="primary" size="small" @click="openClusterEditor(scope.row)">Edit</el-button>
        <el-button link type="primary" size="small" @click="startCluster(scope.row.name)">Start</el-button>
        <el-button link type="danger" size="small" @click="stopCluster(scope.row.name)">Stop</el-button>
        <el-button link type="danger" size="small" @click="deleteCluster(scope.row)">Delete</el-button>
      </template>
    </el-table-column>
  </el-table>

  <el-dialog v-model="creationDialogVisible" title="Create Cluster">
    <el-form :model="clusterForm">
      <el-form-item label="Name">
        <el-input v-model="clusterForm.name" autocomplete="off" />
      </el-form-item>
      <el-form-item label="Servers">
        <el-input-number v-model="clusterForm.servers" :min="1" :max="5" />
      </el-form-item>
      <el-form-item label="Agents">
        <el-input-number v-model="clusterForm.agents" :min="1" :max="5" />
      </el-form-item>
      <el-form-item label="Port">
        <el-input-number v-model="clusterForm.port" :min="30000" :max="31000" />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="openClusterCreationDialog = false">Cancel</el-button>
        <el-button type="primary" @click="createCluster">
          Confirm
        </el-button>
      </span>
    </template>
  </el-dialog>

  <el-drawer
    v-model="clusterDetail"
    :title="currentCluster.name"
    size="80%"
  >
    <div>kubeconfig</div>
    <Codemirror v-model="currentCluster.kubeconfig"/>

    <el-table :data="currentCluster.nodes" stripe style="width: 100%">
      <el-table-column prop="name" label="Name"/>
      <el-table-column prop="role" label="Role"/>
      <el-table-column prop="status" label="Status"/>
    </el-table>
  </el-drawer>

  <el-drawer
    v-model="helmInstallDialog"
    :title="currentCluster.name"
    size="80%"
  >
    <el-form :model="helmForm">
        <el-form-item label="Namespace">
          <el-input v-model="helmForm.namespace" />
        </el-form-item>
        <el-form-item label="Name">
          <el-input v-model="helmForm.name" />
        </el-form-item>
        <el-form-item label="URL">
          <el-input v-model="helmForm.repoURL" />
        </el-form-item>
        <el-form-item label="Version">
          <el-input v-model="helmForm.version" />
        </el-form-item>
        <el-form-item label="Values">
          <el-input v-model="helmForm.valuesStr" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="installHelmChart">
            Confirm
          </el-button>
        </span>
      </template>
  </el-drawer>
</template>
