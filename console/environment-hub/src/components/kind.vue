<script lang="ts" setup>
import { ref, reactive } from 'vue'
import { GetKindClusters, CreateKindClusters, DeleteKindCluster } from './api'

interface Cluster {
  name: string
  servers: number
  agents: number
  port: number
}

const tableData = ref([])

function loadClusters() {
  GetKindClusters((d) => {
    tableData.value = d
  })
}
loadClusters()

function openClusterCreationDialog() {
  creationDialogVisible.value = true
}

function createCluster() {
  creationDialogVisible.value = false
  CreateKindClusters(clusterCreationForm, () => {
    loadClusters()
  })
}

function deleteCluster(row: Cluster) {
  DeleteKindCluster(row.name, (d) => {
    loadClusters()
  })
}

const creationDialogVisible = ref(false)
const clusterCreationForm = reactive({
  servers: 1,
  agents: 1,
  port: 30000
} as Cluster)
</script>

<template>
  <el-button @click="openClusterCreationDialog">Create</el-button>

  <el-table :data="tableData" stripe style="width: 100%">
    <el-table-column prop="name" label="Name" width="180" />
    <el-table-column prop="nodes" label="Nodes" />
    <el-table-column prop="portBinding" label="Port" />
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
</template>
