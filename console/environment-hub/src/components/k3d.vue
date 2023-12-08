<script lang="ts" setup>
import { ref, reactive } from 'vue'
import { GetClusters, GetCluster, CreateClusters, DeleteCluster } from './api'
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
  CreateClusters(clusterCreationForm, () => {
    loadClusters()
  })
}

function deleteCluster(row: Cluster) {
  DeleteCluster(row.name, (d) => {
    loadClusters()
  })
}

const creationDialogVisible = ref(false)
const clusterCreationForm = reactive({
  servers: 1,
  agents: 1,
  port: 30000
} as Cluster)

const clusterDetail = ref(false)
const currentCluster = ref({} as Cluster)
function showClusterDetail(name: string) {
  clusterDetail.value = true
  GetCluster(name, (d) => {
    currentCluster.value = d
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
        <el-button link type="danger" size="small" @click="deleteCluster(scope.row)">Delete</el-button>
      </template>
    </el-table-column>
  </el-table>

  <el-dialog v-model="creationDialogVisible" title="Create Cluster">
    <el-form :model="clusterCreationForm">
      <el-form-item label="Name">
        <el-input v-model="clusterCreationForm.name" autocomplete="off" />
      </el-form-item>
      <el-form-item label="Servers">
        <el-input-number v-model="clusterCreationForm.servers" :min="1" :max="5" />
      </el-form-item>
      <el-form-item label="Agents">
        <el-input-number v-model="clusterCreationForm.agents" :min="1" :max="5" />
      </el-form-item>
      <el-form-item label="Port">
        <el-input-number v-model="clusterCreationForm.port" :min="30000" :max="31000" />
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
</template>
