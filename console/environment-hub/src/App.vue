
<script lang="ts" setup>
import {
  Document,
  Menu as IconMenu,
  Location,
  Setting,
} from '@element-plus/icons-vue'
import k3d from './components/k3d.vue'
import kind from './components/kind.vue'
import { ref } from 'vue'

const panelName = ref('k3d')
const menuSelect = (key: string) => {
  panelName.value = key
}
</script>

<template>
  <div class="common-layout" style="width: 100%; height: 100%; padding: 5px; left: 0px; top: 0px; position: absolute;">
    <el-container>
      <el-aside width="200px">
        <el-row class="tac">
          <el-col :span="12">
            <h5 class="mb-2"></h5>
            <el-menu
              :default-active="panelName"
              class="el-menu-vertical-demo"
              @select="menuSelect"
            >
              <el-menu-item index="k3d">
                <el-icon><icon-menu /></el-icon>
                <span>K3D</span>
              </el-menu-item>
              <el-menu-item index="kind">
                <el-icon><Location /></el-icon>
                <span>Kind</span>
              </el-menu-item>
            </el-menu>
          </el-col>
        </el-row>
      </el-aside>

      <el-main>
        <k3d v-if="panelName === 'k3d'"/>
        <kind v-if="panelName === 'kind'"/>
      </el-main>
    </el-container>
  </div>
</template>
