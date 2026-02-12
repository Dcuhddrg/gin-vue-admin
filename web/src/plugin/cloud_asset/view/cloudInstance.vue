<template>
  <div class="cloud-instance-container">
    <!-- 统计卡片区域 -->
    <div v-if="stats" class="stats-grid">
      <div class="stat-card stat-total">
        <div class="stat-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="2" y="3" width="20" height="14" rx="2" />
            <line x1="8" y1="21" x2="16" y2="21" />
            <line x1="12" y1="17" x2="12" y2="21" />
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.totalCount }}</div>
          <div class="stat-label">实例总数</div>
        </div>
      </div>
      <div class="stat-card stat-running">
        <div class="stat-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2" />
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.runningCount }}</div>
          <div class="stat-label">运行中</div>
        </div>
      </div>
      <div class="stat-card stat-stopped">
        <div class="stat-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="6" y="4" width="4" height="16" />
            <rect x="14" y="4" width="4" height="16" />
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.stoppedCount }}</div>
          <div class="stat-label">已停止</div>
        </div>
      </div>
    </div>

    <!-- 搜索区域 -->
    <div class="search-section">
      <el-form ref="searchForm" :inline="true" :model="searchInfo" class="search-form">
        <el-form-item label="关联项目">
          <el-select
            v-model="searchInfo.projectId"
            placeholder="请选择项目"
            clearable
            filterable
            class="search-select"
          >
            <el-option
              v-for="item in projectList"
              :key="item.ID"
              :label="item.name"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="云厂商类型">
          <el-select
            v-model="searchInfo.providerType"
            placeholder="请选择类型"
            clearable
            class="search-select"
          >
            <el-option
              v-for="item in PROVIDER_OPTIONS"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-select
            v-model="searchInfo.region"
            placeholder="请选择区域"
            clearable
            filterable
            class="search-select"
            :loading="configLoading"
            :disabled="!searchInfo.projectId || !searchInfo.providerType"
          >
            <el-option
              v-for="region in regionList"
              :key="region"
              :label="region"
              :value="region"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="实例状态">
          <el-select v-model="searchInfo.status" placeholder="请选择状态" clearable class="search-select">
            <el-option label="运行中" value="Running" />
            <el-option label="已停止" value="Stopped" />
            <el-option label="启动中" value="Starting" />
            <el-option label="停止中" value="Stopping" />
          </el-select>
        </el-form-item>
        <el-form-item label="实例名称">
          <el-input
            v-model="searchInfo.instanceName"
            placeholder="请输入实例名称"
            clearable
            class="search-input"
          />
        </el-form-item>
        <el-form-item label="实例ID">
          <el-input
            v-model="searchInfo.instanceId"
            placeholder="请输入实例ID"
            clearable
            class="search-input"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" class="search-btn" @click="onSubmit">
            <svg class="btn-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8" />
              <line x1="21" y1="21" x2="16.65" y2="16.65" />
            </svg>
            查询
          </el-button>
          <el-button class="reset-btn" @click="onReset">
            <svg class="btn-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
              <path d="M3 3v5h5" />
            </svg>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 操作按钮区域 -->
    <div class="action-bar">
      <div class="left-actions">
        <el-button type="primary" class="sync-btn" @click="openSyncDialog">
          <svg class="btn-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12a9 9 0 0 1-9 9m9-9a9 9 0 0 0-9-9m9 9H3m9 9a9 9 0 0 1-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 0 1 9-9" />
          </svg>
          同步实例
        </el-button>
        <el-button class="refresh-btn" @click="refreshTable">
          <svg class="btn-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2" />
          </svg>
          刷新
        </el-button>
        <el-button class="clear-cache-btn" @click="clearCache">
          <svg class="btn-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 6h18M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
          </svg>
          清除缓存
        </el-button>
      </div>
      <div class="right-actions">
        <el-tag v-if="lastSyncTime" class="sync-time-tag" type="info">
          上次同步: {{ formatLastSyncTime(lastSyncTime) }}
        </el-tag>
      </div>
    </div>

    <!-- 表格区域 -->
    <div class="table-wrapper">
      <el-table
        v-loading="loading"
        :data="tableData"
        row-key="ID"
        class="instance-table"
        :header-cell-class-name="'table-header-cell'"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="实例信息" min-width="280">
          <template #default="scope">
            <div class="instance-info">
              <div class="instance-name">{{ scope.row.instanceName || scope.row.instanceId }}</div>
              <div class="instance-id">{{ scope.row.instanceId }}</div>
              <el-tag v-if="scope.row.instanceType" size="small" class="instance-type-tag">
                {{ scope.row.instanceType }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="云厂商" width="120">
          <template #default="scope">
            <div class="provider-cell">
              <el-tag size="small" :class="`provider-tag provider-${scope.row.providerType}`">
                {{ getProviderLabel(scope.row.providerType) }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="区域" width="120" prop="region">
          <template #default="scope">
            <el-tag type="info" size="small">{{ scope.row.region }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="配置" width="180">
          <template #default="scope">
            <div class="config-cell">
              <div v-if="scope.row.cpu" class="config-item">
                <svg class="config-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="4" y="4" width="16" height="16" rx="2" />
                  <rect x="9" y="9" width="6" height="6" />
                </svg>
                <span>{{ scope.row.cpu }}核</span>
              </div>
              <div v-if="scope.row.memory" class="config-item">
                <svg class="config-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M2 12h20M2 12a10 10 0 0 1 10-10v20a10 10 0 0 1-10-10" />
                  <path d="M12 2v20" />
                </svg>
                <span>{{ scope.row.memory }}GB</span>
              </div>
              <div v-if="scope.row.diskSize" class="config-item">
                <svg class="config-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10" />
                  <circle cx="12" cy="12" r="3" />
                </svg>
                <span>{{ scope.row.diskSize }}GB</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="IP地址" width="220">
          <template #default="scope">
            <div class="ip-cell">
              <div v-if="scope.row.publicIp" class="ip-item">
                <svg class="ip-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10" />
                  <line x1="2" y1="12" x2="22" y2="12" />
                  <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z" />
                </svg>
                <span class="ip-label">公网:</span>
                <span class="ip-value">{{ scope.row.publicIp }}</span>
                <el-button
                  type="primary"
                  link
                  size="small"
                  class="copy-btn"
                  @click="copyText(scope.row.publicIp)"
                >
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" />
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
                  </svg>
                </el-button>
              </div>
              <div v-if="scope.row.privateIp" class="ip-item ip-private">
                <svg class="ip-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="11" width="18" height="10" rx="2" />
                  <circle cx="12" cy="5" r="2" />
                  <path d="M12 7v4" />
                  <line x1="8" y1="16" x2="8" y2="16" />
                  <line x1="16" y1="16" x2="16" y2="16" />
                </svg>
                <span class="ip-label">内网:</span>
                <span class="ip-value">{{ scope.row.privateIp }}</span>
                <el-button
                  type="primary"
                  link
                  size="small"
                  class="copy-btn"
                  @click="copyText(scope.row.privateIp)"
                >
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" />
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
                  </svg>
                </el-button>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120">
          <template #default="scope">
            <div class="status-cell">
              <div :class="`status-indicator status-${getStatusClass(scope.row.status)}`"></div>
              <span class="status-text">{{ getStatusLabel(scope.row.status) }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="160">
          <template #default="scope">
            <div v-if="scope.row.createdTime" class="time-cell">
              {{ formatDateTime(scope.row.createdTime) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="最后同步" width="160">
          <template #default="scope">
            <div v-if="scope.row.lastSyncAt" class="time-cell sync-time">
              {{ formatDateTime(scope.row.lastSyncAt) }}
            </div>
            <div v-else class="time-cell sync-time not-synced">
              未同步
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="scope">
            <div class="action-buttons">
              <el-button
                type="primary"
                link
                size="small"
                class="action-btn view-btn"
                @click="viewDetail(scope.row)"
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
                  <circle cx="12" cy="12" r="3" />
                </svg>
                详情
              </el-button>
              <el-button
                type="danger"
                link
                size="small"
                class="action-btn delete-btn"
                @click="deleteInstance(scope.row)"
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6" />
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
                </svg>
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 同步对话框 -->
    <el-dialog
      v-model="syncDialogVisible"
      title="同步云服务器实例"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="syncForm" label-width="100px">
        <el-form-item label="关联项目" required>
          <el-select
            v-model="syncForm.projectId"
            placeholder="请选择项目"
            filterable
            class="full-width"
          >
            <el-option
              v-for="item in projectList"
              :key="item.ID"
              :label="item.name"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="云厂商类型" required>
          <el-select
            v-model="syncForm.providerType"
            placeholder="请选择类型"
            filterable
            class="full-width"
            :disabled="!syncForm.projectId"
          >
            <el-option
              v-for="item in PROVIDER_OPTIONS"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="false" label="同步方式">
          <el-radio-group v-model="syncForm.syncType">
            <el-radio value="single">单个区域同步</el-radio>
            <el-radio value="batch">批量区域同步</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="syncForm.syncType === 'single'" label="区域">
          <el-select
            v-model="syncForm.region"
            placeholder="请选择区域"
            filterable
            class="full-width"
            disabled
            :loading="configLoading"
          >
            <el-option
              v-for="region in providerRegions"
              :key="region"
              :label="region"
              :value="region"
            />
          </el-select>
          <div v-if="syncForm.region" class="form-tip">
            系统自动检测到唯一可用区域，已准备就绪。
          </div>
        </el-form-item>
        <el-form-item v-if="syncForm.syncType === 'batch'" label="区域">
          <el-select
            v-model="syncForm.regions"
            placeholder="请选择要同步的区域"
            multiple
            filterable
            class="full-width"
            disabled
            :loading="configLoading"
          >
            <el-option
              v-for="region in providerRegions"
              :key="region"
              :label="region"
              :value="region"
            />
          </el-select>
          <div v-if="syncForm.regions.length > 0" class="form-tip">
            系统自动检测到 {{ syncForm.regions.length }} 个可用区域，已自动全选。
          </div>
        </el-form-item>
        <el-form-item label="强制同步">
          <el-switch v-model="syncForm.forceSync" />
          <div class="form-tip">开启后将忽略缓存，直接从云厂商获取最新数据</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="syncDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="syncing" @click="handleSync">
          开始同步
        </el-button>
      </template>
    </el-dialog>

    <!-- 实例详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="`实例详情 - ${currentInstance?.instanceName || currentInstance?.instanceId}`"
      width="700px"
    >
      <div v-if="currentInstance" class="instance-detail">
        <div class="detail-section">
          <div class="detail-title">基础信息</div>
          <div class="detail-grid">
            <div class="detail-item">
              <span class="detail-label">实例ID:</span>
              <span class="detail-value">{{ currentInstance.instanceId }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">实例名称:</span>
              <span class="detail-value">{{ currentInstance.instanceName || '-' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">状态:</span>
              <span :class="`detail-value status-${getStatusClass(currentInstance.status)}`">
                {{ getStatusLabel(currentInstance.status) }}
              </span>
            </div>
            <div class="detail-item">
              <span class="detail-label">区域:</span>
              <span class="detail-value">{{ currentInstance.region }}</span>
            </div>
          </div>
        </div>
        <div class="detail-section">
          <div class="detail-title">配置信息</div>
          <div class="detail-grid">
            <div class="detail-item">
              <span class="detail-label">实例规格:</span>
              <span class="detail-value">{{ currentInstance.instanceType || '-' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">CPU:</span>
              <span class="detail-value">{{ currentInstance.cpu ? `${currentInstance.cpu}核` : '-' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">内存:</span>
              <span class="detail-value">{{ currentInstance.memory ? `${currentInstance.memory}GB` : '-' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">磁盘:</span>
              <span class="detail-value">{{ currentInstance.diskSize ? `${currentInstance.diskSize}GB` : '-' }}</span>
            </div>
            <div class="detail-item full-width">
              <span class="detail-label">操作系统:</span>
              <span class="detail-value">
                {{ currentInstance.osName || '-' }}
                <span v-if="currentInstance.osVersion" class="os-version">
                  ({{ currentInstance.osVersion }})
                </span>
              </span>
            </div>
          </div>
        </div>
        <div class="detail-section">
          <div class="detail-title">网络信息</div>
          <div class="detail-grid">
            <div class="detail-item full-width">
              <span class="detail-label">公网IP:</span>
              <span class="detail-value">{{ currentInstance.publicIp || '-' }}</span>
            </div>
            <div class="detail-item full-width">
              <span class="detail-label">内网IP:</span>
              <span class="detail-value">{{ currentInstance.privateIp || '-' }}</span>
            </div>
          </div>
        </div>
        <div class="detail-section">
          <div class="detail-title">时间信息</div>
          <div class="detail-grid">
            <div class="detail-item">
              <span class="detail-label">创建时间:</span>
              <span class="detail-value">
                {{ currentInstance.createdTime ? formatDateTime(currentInstance.createdTime) : '-' }}
              </span>
            </div>
            <div class="detail-item">
              <span class="detail-label">到期时间:</span>
              <span class="detail-value">
                {{ currentInstance.expiredTime ? formatDateTime(currentInstance.expiredTime) : '-' }}
              </span>
            </div>
            <div class="detail-item full-width">
              <span class="detail-label">付费类型:</span>
              <span class="detail-value">{{ getChargeTypeLabel(currentInstance.chargeType) }}</span>
            </div>
          </div>
        </div>
        <div v-if="currentInstance.remark" class="detail-section">
          <div class="detail-title">备注</div>
          <div class="detail-content">{{ currentInstance.remark }}</div>
        </div>
      </div>
      <template #footer>
        <el-button type="primary" @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  getCloudInstanceVOList,
  getCloudInstanceById,
  deleteCloudInstance,
  batchDeleteCloudInstance,
  syncInstances,
  batchSyncInstances,
  getInstanceStats,
  clearInstanceCache
} from '@/plugin/cloud_asset/api/cloudInstance.js'
import { getProjectList } from '@/plugin/project_manager/api/project.js'
import { PROVIDER_OPTIONS, getProviderLabel } from '@/plugin/cloud_asset/config/provider.js'

import { getCloudProviderList, getProviderConfig } from '@/plugin/cloud_asset/api/cloudProvider.js'

  import { ref, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatTimeToStr } from '@/utils/date'

defineOptions({
  name: 'CloudInstance'
})

// 数据状态
const loading = ref(false)
const syncing = ref(false)
const configLoading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const tableData = ref([])
const providerList = ref([])
const projectList = ref([])
const regionList = ref([])
const stats = ref(null)
const lastSyncTime = ref(null)
const providerConfig = ref(null)

// 搜索条件
const searchInfo = ref({
  projectId: undefined,
  providerType: '',
  region: '',
  status: '',
  instanceName: '',
  instanceId: ''
})

// 同步表单
const syncDialogVisible = ref(false)
const syncForm = ref({
  projectId: undefined,
  providerType: '',
  providerId: undefined,
  region: '',
  regions: [],
  syncType: 'single',
  forceSync: false
})
const providerRegions = ref([])

// 实例详情
const detailDialogVisible = ref(false)
const currentInstance = ref(null)

// 缓存
const configCache = new Map()

const getProjects = async () => {
  const res = await getProjectList({ page: 1, pageSize: 1000 })
  if (res.code === 0) {
    projectList.value = res.data.list
  }
}

// 获取云厂商列表
const getProviders = async () => {
  const res = await getCloudProviderList({ page: 1, pageSize: 1000 })
  if (res.code === 0) {
    providerList.value = res.data.list
  }
}

// 获取活跃的云厂商列表
const activeProviderList = ref([])
const getActiveProviders = async () => {
  const res = await getCloudProviderList({ page: 1, pageSize: 1000, status: 1 })
  if (res.code === 0) {
    activeProviderList.value = res.data.list
  }
}

// 获取厂商配置（核心联动逻辑）
const fetchProviderConfig = async (projectId, providerType, target = 'search') => {
  if (!projectId || !providerType) return

  const cacheKey = `${projectId}_${providerType}`
  if (configCache.has(cacheKey)) {
    handleConfigSuccess(configCache.get(cacheKey), target)
    return
  }

  configLoading.value = true
  try {
    const res = await getProviderConfig({ projectId, providerType })
    if (res.code === 0) {
      configCache.set(cacheKey, res.data)
      handleConfigSuccess(res.data, target)
    }
  } catch (error) {
    ElMessage.error('获取厂商配置失败，请检查凭证配置')
    handleConfigError(target)
  } finally {
    configLoading.value = false
  }
}

const handleConfigSuccess = (data, target) => {
  const regions = data.regions.map(r => r.regionId)
  if (target === 'search') {
    regionList.value = regions
  } else if (target === 'sync') {
    providerRegions.value = regions
    // 自动判断同步方式
    if (regions.length > 1) {
      // 多区域 -> 批量同步
      syncForm.value.syncType = 'batch'
      syncForm.value.regions = regions // 全选
      syncForm.value.region = '' // 清空单选字段
    } else if (regions.length === 1) {
      // 单区域 -> 单个同步
      syncForm.value.syncType = 'single'
      syncForm.value.region = regions[0] // 选中唯一区域
      syncForm.value.regions = [] // 清空多选字段
    } else {
      // 无区域
      syncForm.value.syncType = 'single'
      syncForm.value.region = ''
      syncForm.value.regions = []
    }
  }
  providerConfig.value = data
}

const handleConfigError = (target) => {
  if (target === 'search') {
    regionList.value = []
    searchInfo.value.region = ''
  } else if (target === 'sync') {
    providerRegions.value = []
    syncForm.value.region = ''
  }
}

// 监听搜索表单联动
watch(
  [() => searchInfo.value.projectId, () => searchInfo.value.providerType],
  ([newPid, newType]) => {
    if (newPid && newType) {
      fetchProviderConfig(newPid, newType, 'search')
    } else {
      regionList.value = []
      searchInfo.value.region = ''
    }
  }
)

// 监听同步表单项目和类型变化，自动设置ProviderID
watch(
  [() => syncForm.value.projectId, () => syncForm.value.providerType],
  ([newPid, newType]) => {
    if (newPid && newType) {
      const provider = activeProviderList.value.find(p => p.projectId === newPid && p.type === newType)
      if (provider) {
        syncForm.value.providerId = provider.ID
      } else {
        syncForm.value.providerId = undefined
        if (newPid && newType) {
          ElMessage.warning('该项目下未配置对应的云厂商凭证')
        }
      }
    } else {
      syncForm.value.providerId = undefined
    }
  }
)

// 监听同步表单联动（需要根据选择的厂商ID反查项目ID和类型）
watch(() => syncForm.value.providerId, (newVal) => {
  if (newVal) {
    const provider = activeProviderList.value.find(p => p.ID === newVal)
    if (provider) {
      fetchProviderConfig(provider.projectId, provider.type, 'sync')
    }
  } else {
    providerRegions.value = []
    syncForm.value.region = ''
  }
})

// 获取统计数据
const getStats = async () => {
  const res = await getInstanceStats()
  if (res.code === 0) {
    stats.value = res.data
  }
}

// 获取表格数据
const getTableData = async () => {
  loading.value = true
  try {
    const res = await getCloudInstanceVOList({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value
    })
    if (res.code === 0) {
      tableData.value = res.data.list
      total.value = res.data.total
      page.value = res.data.page
      pageSize.value = res.data.pageSize
      // 更新最后同步时间
      if (res.data.list.length > 0) {
        const syncTimes = res.data.list
          .map(item => item.lastSyncAt)
          .filter(time => time)
          .sort((a, b) => new Date(b) - new Date(a))
        if (syncTimes.length > 0) {
          lastSyncTime.value = syncTimes[0]
        }
      }
    }
  } finally {
    loading.value = false
  }
}

// 获取云厂商的区域列表
const getProviderRegions = (providerId) => {
  const provider = providerList.value.find(p => p.ID === providerId)
  if (provider && provider.region) {
    return [provider.region]
  }
  return regionList.value
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 搜索
const onSubmit = () => {
  page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = {
    projectId: undefined,
    providerType: '',
    region: '',
    status: '',
    instanceName: '',
    instanceId: ''
  }
  page.value = 1
  getTableData()
}

// 刷新表格
const refreshTable = () => {
  getTableData()
  getStats()
}

// 清除缓存
const clearCache = async () => {
  try {
    await ElMessageBox.confirm('确定要清除云服务器实例缓存吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await clearInstanceCache({
      providerId: searchInfo.value.providerId,
      region: searchInfo.value.region
    })
    if (res.code === 0) {
      ElMessage.success('缓存清除成功')
      refreshTable()
    }
  } catch (error) {
    // 用户取消
  }
}

// 打开同步对话框
const openSyncDialog = () => {
  syncForm.value = {
    projectId: undefined,
    providerType: '',
    providerId: undefined,
    region: '',
    regions: [],
    syncType: 'single',
    forceSync: false
  }
  syncDialogVisible.value = true
}

// 同步实例
const handleSync = async () => {
  if (!syncForm.value.providerId) {
    ElMessage.warning('请选择有效的项目和云厂商类型')
    return
  }

  syncing.value = true
  try {
    let res
    if (syncForm.value.syncType === 'single') {
      if (!syncForm.value.region) {
        ElMessage.warning('请选择区域')
        syncing.value = false
        return
      }
      res = await syncInstances({
        providerId: syncForm.value.providerId,
        region: syncForm.value.region,
        forceSync: syncForm.value.forceSync
      })
    } else {
      if (!syncForm.value.regions || syncForm.value.regions.length === 0) {
        ElMessage.warning('请选择要同步的区域')
        syncing.value = false
        return
      }
      res = await batchSyncInstances({
        providerId: syncForm.value.providerId,
        regions: syncForm.value.regions,
        forceSync: syncForm.value.forceSync
      })
    }

    if (res.code === 0) {
      ElMessage.success('同步成功')
      syncDialogVisible.value = false
      refreshTable()
    }
  } finally {
    syncing.value = false
  }
}

// 查看详情
const viewDetail = async (row) => {
  const res = await getCloudInstanceById({ id: row.ID })
  if (res.code === 0) {
    currentInstance.value = res.data
    detailDialogVisible.value = true
  }
}

// 删除实例
const deleteInstance = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要删除实例 "${row.instanceName || row.instanceId}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await deleteCloudInstance({ id: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      refreshTable()
    }
  } catch (error) {
    // 用户取消
  }
}

// 工具函数
const getStatusClass = (status) => {
  const map = {
    Running: 'running',
    Stopped: 'stopped',
    Starting: 'starting',
    Stopping: 'stopping'
  }
  return map[status] || 'unknown'
}

const getStatusLabel = (status) => {
  const labels = {
    Running: '运行中',
    Stopped: '已停止',
    Starting: '启动中',
    Stopping: '停止中'
  }
  return labels[status] || status || '未知'
}

const getChargeTypeLabel = (type) => {
  const labels = {
    'PostPaid': '按量付费',
    'PrePaid': '包年包月',
    'SpotWithPriceLimit': '竞价实例'
  }
  return labels[type] || type || '-'
}

const formatDateTime = (time) => {
  if (!time) return '-'
  return formatTimeToStr(time, 'yyyy-MM-dd hh:mm:ss')
}

const formatLastSyncTime = (time) => {
  if (!time) return '未同步'
  const now = new Date()
  const diff = now - new Date(time)
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (days > 0) return `${days}天前`
  if (hours > 0) return `${hours}小时前`
  if (minutes > 0) return `${minutes}分钟前`
  return '刚刚'
}

const copyText = (text) => {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success('复制成功')
  })
}

// 初始化
onMounted(() => {
  getProjects()
  getProviders()
  getActiveProviders()
  getTableData()
  getStats()
})

// 监听云厂商选择变化
watch(() => syncForm.value.providerId, (newVal) => {
  if (newVal) {
    providerRegions.value = getProviderRegions(newVal)
    // 设置默认区域
    if (providerRegions.value.length === 1) {
      syncForm.value.region = providerRegions.value[0]
    }
  } else {
    providerRegions.value = []
  }
})
</script>

<style lang="scss" scoped>
.cloud-instance-container {
  padding: 24px;
  background: #f8fafc;
  min-height: 100vh;
}

// 统计卡片
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
  border-radius: 12px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
  border: 1px solid rgba(0, 0, 0, 0.04);

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  }

  .stat-icon {
    width: 56px;
    height: 56px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
  }

  .stat-content {
    flex: 1;
  }

  .stat-value {
    font-size: 32px;
    font-weight: 700;
    color: #1e293b;
    line-height: 1;
    margin-bottom: 4px;
  }

  .stat-label {
    font-size: 14px;
    color: #64748b;
  }
}

.stat-total .stat-icon {
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
}

.stat-running .stat-icon {
  background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
}

.stat-stopped .stat-icon {
  background: linear-gradient(135deg, #f97316 0%, #ea580c 100%);
}

// 搜索区域
.search-section {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.search-form {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;

  :deep(.el-form-item) {
    margin-bottom: 0;
  }
}

.search-select,
.search-input {
  width: 180px;
}

.search-btn,
.reset-btn {
  height: 32px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  gap: 6px;

  .btn-icon {
    width: 16px;
    height: 16px;
  }
}

.search-btn {
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  border: none;

  &:hover {
    background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%);
  }
}

// 操作栏
.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 16px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.left-actions {
  display: flex;
  gap: 12px;
}

.sync-btn,
.refresh-btn,
.clear-cache-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 16px;
  height: 36px;

  .btn-icon {
    width: 18px;
    height: 18px;
  }
}

.sync-btn {
  background: linear-gradient(135deg, #8b5cf6 0%, #6d28d9 100%);
  border: none;

  &:hover {
    background: linear-gradient(135deg, #7c3aed 0%, #5b21b6 100%);
  }
}

.sync-time-tag {
  font-size: 13px;
}

// 表格区域
.table-wrapper {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  overflow: hidden;
}

.instance-table {
  :deep(.el-table__header-wrapper) {
    background: #f1f5f9;
  }

  :deep(.table-header-cell) {
    background: #f1f5f9 !important;
    color: #475569;
    font-weight: 600;
  }

  :deep(.el-table__row) {
    transition: all 0.2s ease;

    &:hover {
      background: #f8fafc !important;
    }
  }
}

.instance-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.instance-name {
  font-weight: 600;
  color: #1e293b;
  font-size: 14px;
}

.instance-id {
  font-size: 12px;
  color: #94a3b8;
  font-family: 'SF Mono', Monaco, monospace;
}

.instance-type-tag {
  align-self: flex-start;
  margin-top: 4px;
  background: #e0e7ff;
  color: #4338ca;
  border: none;
}

.provider-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.provider-name {
  font-size: 12px;
  color: #94a3b8;
}

.provider-tag {
  border: none;
  font-weight: 500;
}

.provider-aliyun { background: linear-gradient(135deg, #ff9a56 0%, #ff6b35 100%); color: white; }
.provider-tencent { background: linear-gradient(135deg, #00a4ff 0%, #0073e6 100%); color: white; }
.provider-aws { background: linear-gradient(135deg, #ff9900 0%, #ff6600 100%); color: white; }
.provider-huawei { background: linear-gradient(135deg, #c40d15 0%, #a30b12 100%); color: white; }
.provider-baidu { background: linear-gradient(135deg, #2932e1 0%, #1f26b5 100%); color: white; }
.provider-default { background: #e5e7eb; color: #374151; }

.config-cell {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.config-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #475569;

  .config-icon {
    width: 14px;
    height: 14px;
    color: #64748b;
  }
}

.ip-cell {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.ip-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
}

.ip-private {
  color: #64748b;
}

.ip-icon {
  width: 14px;
  height: 14px;
  color: #64748b;
}

.ip-label {
  color: #94a3b8;
  font-size: 11px;
}

.ip-value {
  font-family: 'SF Mono', Monaco, monospace;
  color: #1e293b;
}

.copy-btn {
  padding: 2px;
  margin-left: auto;
  color: #64748b;

  &:hover {
    color: #3b82f6;
  }

  svg {
    width: 12px;
    height: 12px;
  }
}

.status-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.status-running {
  background: #22c55e;
  box-shadow: 0 0 8px rgba(34, 197, 94, 0.4);
}

.status-stopped {
  background: #ef4444;
  box-shadow: 0 0 8px rgba(239, 68, 68, 0.4);
  animation: none;
}

.status-starting {
  background: #eab308;
  box-shadow: 0 0 8px rgba(234, 179, 8, 0.4);
}

.status-stopping {
  background: #f97316;
  box-shadow: 0 0 8px rgba(249, 115, 22, 0.4);
}

.status-unknown {
  background: #94a3b8;
  animation: none;
}

.status-text {
  font-size: 13px;
  font-weight: 500;
}

.status-running { color: #16a34a; }
.status-stopped { color: #dc2626; }
.status-starting { color: #ca8a04; }
.status-stopping { color: #ea580c; }
.status-unknown { color: #64748b; }

.time-cell {
  font-size: 13px;
  color: #475569;
  font-family: 'SF Mono', Monaco, monospace;
}

.sync-time {
  color: #64748b;
  font-style: italic;
}

.sync-time.not-synced {
  color: #94a3b8;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;

  svg {
    width: 14px;
    height: 14px;
  }
}

.delete-btn {
  color: #ef4444;

  &:hover {
    color: #dc2626;
  }
}

// 分页
.pagination-wrapper {
  padding: 16px 20px;
  display: flex;
  justify-content: flex-end;
}

// 对话框
.full-width {
  width: 100%;
}

.form-tip {
  font-size: 12px;
  color: #64748b;
  margin-top: 4px;
}

.provider-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.provider-tag.provider-aliyun,
.provider-tag.provider-tencent,
.provider-tag.provider-aws,
.provider-tag.provider-huawei,
.provider-tag.provider-baidu {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

// 实例详情
.instance-detail {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.detail-section {
  background: #f8fafc;
  border-radius: 8px;
  padding: 16px;
}

.detail-title {
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 2px solid #e2e8f0;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;

  &.full-width {
    grid-column: 1 / -1;
  }
}

.detail-label {
  font-size: 12px;
  color: #64748b;
}

.detail-value {
  font-size: 14px;
  color: #1e293b;
  font-weight: 500;

  &.status-running { color: #16a34a; }
  &.status-stopped { color: #dc2626; }
  &.status-starting { color: #ca8a04; }
  &.status-stopping { color: #ea580c; }
}

.os-version {
  color: #64748b;
  font-weight: 400;
}

.detail-content {
  font-size: 14px;
  color: #475569;
  line-height: 1.6;
  white-space: pre-wrap;
}
</style>
