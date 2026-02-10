<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo" class="flex flex-wrap gap-4">
        <el-form-item label="关联项目">
          <el-select v-model="searchInfo.projectId" placeholder="请选择关联项目" clearable filterable>
            <el-option
              v-for="item in projectList"
              :key="item.ID"
              :label="item.name"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="厂商类型">
          <el-select v-model="searchInfo.type" placeholder="请选择厂商类型" clearable>
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="华为云" value="huawei" />
            <el-option label="百度云" value="baidu" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">
            查询
          </el-button>
          <el-button icon="refresh" @click="onReset"> 重置 </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="addCloudProvider">
          新增云厂商
        </el-button>
      </div>
      <el-table :data="tableData" row-key="ID">
        <el-table-column
          align="left"
          label="关联项目"
          min-width="150"
        >
          <template #default="scope">
            <el-tag v-if="scope.row.project" type="info">{{ scope.row.project.name }}</el-tag>
            <span v-else>未关联</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="ID" min-width="80" prop="ID" />
        <el-table-column
          align="left"
          label="厂商类型"
          min-width="120"
          prop="type"
        >
          <template #default="scope">
            <el-tag v-if="scope.row.type === 'aliyun'" type="success">阿里云</el-tag>
            <el-tag v-else-if="scope.row.type === 'tencent'" type="primary">腾讯云</el-tag>
            <el-tag v-else-if="scope.row.type === 'aws'" type="warning">AWS</el-tag>
            <el-tag v-else-if="scope.row.type === 'huawei'" type="danger">华为云</el-tag>
            <el-tag v-else-if="scope.row.type === 'baidu'" type="info">百度云</el-tag>
            <el-tag v-else>{{ scope.row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="AccessKey"
          min-width="200"
          prop="ak"
        >
          <template #default="scope">
            <span class="text-gray-500">{{ maskSensitiveInfo(scope.row.ak) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="SecretKey"
          min-width="200"
          prop="sk"
        >
          <template #default="scope">
            <span class="text-gray-500">{{ maskSensitiveInfo(scope.row.sk) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="区域"
          min-width="120"
          prop="region"
        />
        <el-table-column align="left" label="状态" min-width="100">
          <template #default="scope">
            <el-switch
              v-model="scope.row.status"
              inline-prompt
              :active-value="1"
              :inactive-value="2"
              @change="switchStatus(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="备注"
          min-width="150"
          prop="remark"
          show-overflow-tooltip
        />
        <el-table-column label="操作" :min-width="appStore.operateMinWith" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteCloudProviderFunc(scope.row)"
            >
              删除
            </el-button>
            <el-button
              type="primary"
              link
              icon="edit"
              @click="openEdit(scope.row)"
            >
              编辑
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-drawer
      v-model="addCloudProviderDialog"
      :size="appStore.drawerSize"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ dialogTitle }}</span>
          <div>
            <el-button @click="closeAddCloudProviderDialog">取 消</el-button>
            <el-button type="primary" @click="enterAddCloudProviderDialog">
              确 定
            </el-button>
          </div>
        </div>
      </template>

      <el-form
        ref="cloudProviderForm"
        :rules="rules"
        :model="cloudProviderInfo"
        label-width="100px"
      >
        <el-form-item label="关联项目" prop="projectId">
          <el-select
            v-model="cloudProviderInfo.projectId"
            placeholder="请选择关联项目"
            class="w-full"
            filterable
          >
            <el-option
              v-for="item in projectList"
              :key="item.ID"
              :label="item.name"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="厂商类型" prop="type">
          <el-select v-model="cloudProviderInfo.type" placeholder="请选择厂商类型" class="w-full">
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="华为云" value="huawei" />
            <el-option label="百度云" value="baidu" />
          </el-select>
        </el-form-item>
        <el-form-item label="AccessKey" prop="ak">
          <el-input
            v-model="cloudProviderInfo.ak"
            placeholder="请输入AccessKey"
            type="password"
            show-password
          />
        </el-form-item>
        <el-form-item label="SecretKey" prop="sk">
          <el-input
            v-model="cloudProviderInfo.sk"
            placeholder="请输入SecretKey"
            type="password"
            show-password
            @blur="handleSkBlur"
          />
        </el-form-item>
        <el-form-item label="区域" prop="region">
          <RegionsSelector
            ref="regionsSelector"
            v-model="cloudProviderInfo.region"
            :provider="cloudProviderInfo.type"
            :access-key="cloudProviderInfo.ak"
            :secret-key="cloudProviderInfo.sk"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch
            v-model="cloudProviderInfo.status"
            inline-prompt
            :active-value="1"
            :inactive-value="2"
            active-text="启用"
            inactive-text="禁用"
          />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="cloudProviderInfo.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息"
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    getCloudProviderList,
    createCloudProvider,
    updateCloudProvider,
    deleteCloudProvider
  } from '@/plugin/cloud_asset/api/cloudProvider.js'
  import { getProjectList as getProjectListApi } from '@/plugin/project_manager/api/project'
  import RegionsSelector from './components/RegionsSelector.vue'

  import { ref, computed, watch } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { useAppStore } from '@/pinia'

  defineOptions({
    name: 'CloudProvider'
  })

  const appStore = useAppStore()

  const searchInfo = ref({
    type: '',
    projectId: undefined
  })

  const onSubmit = () => {
    page.value = 1
    getTableData()
  }

  const onReset = () => {
    searchInfo.value = {
      type: '',
      projectId: undefined
    }
    getTableData()
  }

  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  const projectList = ref([])

  const getProjectList = async () => {
    const res = await getProjectListApi({ page: 1, pageSize: 999 })
    if (res.code === 0) {
      projectList.value = res.data.list
    }
  }

  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  const getTableData = async () => {
    const table = await getCloudProviderList({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value
    })
    if (table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }

  const initPage = async () => {
    getProjectList()
    getTableData()
  }

  initPage()

  const maskSensitiveInfo = (value) => {
    if (!value) return ''
    if (value.length <= 8) {
      return '********'
    }
    return value.substring(0, 4) + '****' + value.substring(value.length - 4)
  }

  const deleteCloudProviderFunc = async (row) => {
    ElMessageBox.confirm('确定要删除该云厂商吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const res = await deleteCloudProvider({ id: row.ID })
      if (res.code === 0) {
        ElMessage.success('删除成功')
        await getTableData()
      }
    })
  }

  const cloudProviderInfo = ref({
    ak: '',
    sk: '',
    type: '',
    region: '',
    status: 1,
    remark: '',
    projectId: undefined
  })

  const rules = ref({
    projectId: [{ required: true, message: '请选择关联项目', trigger: 'change' }],
    type: [{ required: true, message: '请选择厂商类型', trigger: 'change' }],
    ak: [{ required: true, message: '请输入AccessKey', trigger: 'blur' }],
    sk: [{ required: true, message: '请输入SecretKey', trigger: 'blur' }],
    region: [{ required: true, message: '请输入区域', trigger: 'blur' }]
  })

  const cloudProviderForm = ref(null)
  const addCloudProviderDialog = ref(false)
  const dialogFlag = ref('add')
  const regionsSelector = ref(null)

  const dialogTitle = computed(() => {
    return dialogFlag.value === 'add' ? '新增云厂商' : '编辑云厂商'
  })

  const handleSkBlur = () => {
    if (cloudProviderInfo.value.type && cloudProviderInfo.value.ak && cloudProviderInfo.value.sk) {
      regionsSelector.value?.fetchRegions()
    }
  }

  const addCloudProvider = () => {
    dialogFlag.value = 'add'
    addCloudProviderDialog.value = true
  }

  const openEdit = async (row) => {
    dialogFlag.value = 'edit'
    cloudProviderInfo.value = JSON.parse(JSON.stringify(row))
    addCloudProviderDialog.value = true
  }

  const closeAddCloudProviderDialog = () => {
    cloudProviderForm.value?.resetFields()
    cloudProviderInfo.value = {
      ak: '',
      sk: '',
      type: '',
      region: '',
      status: 1,
      remark: '',
      projectId: undefined
    }
    addCloudProviderDialog.value = false
  }

  const enterAddCloudProviderDialog = async () => {
    cloudProviderForm.value.validate(async (valid) => {
      if (valid) {
        const req = {
          ...cloudProviderInfo.value
        }
        let res
        if (dialogFlag.value === 'add') {
          res = await createCloudProvider(req)
          if (res.code === 0) {
            ElMessage({ type: 'success', message: '创建成功' })
            await getTableData()
            closeAddCloudProviderDialog()
          }
        } else {
          res = await updateCloudProvider(req)
          if (res.code === 0) {
            ElMessage({ type: 'success', message: '编辑成功' })
            await getTableData()
            closeAddCloudProviderDialog()
          }
        }
      }
    })
  }

  const switchStatus = async (row) => {
    const req = {
      ...row
    }
    const res = await updateCloudProvider(req)
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: `${req.status === 2 ? '禁用' : '启用'}成功`
      })
      await getTableData()
    } else {
      row.status = req.status === 1 ? 2 : 1
    }
  }
</script>

<style lang="scss" scoped>
</style>
