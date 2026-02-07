<template>
  <div class="w-full">
    <!-- 提示条 -->
    <div class="h-12px text-12px leading-12px text-gray-500 mb-2 flex items-center">
      <el-icon class="mr-1"><InfoFilled /></el-icon>
      AccessKey 与 SecretKey 仅做一次性验证，不会被持久化存储
    </div>

    <!-- 区域选择容器 -->
    <div
      v-loading="regionState.loading"
      element-loading-text="正在获取可用区..."
      class="border rounded p-4 min-h-100px relative flex flex-col gap-4"
      :class="{'border-red-500': regionState.errorMsg}"
    >
      <!-- 顶部操作栏 -->
      <div class="flex justify-end mb-2">
        <el-button
          type="primary"
          link
          size="small"
          :disabled="!canFetchRegions"
          @click="fetchRegions"
        >
          <el-icon class="mr-1"><Refresh /></el-icon>
          验证并获取可用区
        </el-button>
      </div>

      <!-- 错误提示 -->
      <div v-if="regionState.errorMsg" class="text-red-500 text-sm">
        {{ regionState.errorMsg }}
      </div>

      <!-- 空状态 -->
      <div v-if="!regionState.loading && regionState.regions.length === 0 && !regionState.errorMsg" class="text-gray-400 text-center py-4">
        请填写 AccessKey/SecretKey 并选择厂商后获取可用区
      </div>

      <!-- 区域列表 -->
      <el-checkbox-group v-model="regionState.selected" @change="handleSelectionChange">
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          <el-checkbox
            v-for="region in regionState.regions"
            :key="region.regionId"
            :label="region.regionId"
            class="!mr-0 w-full"
          >
            <div class="flex items-center w-full overflow-hidden">
              <span class="truncate" :title="region.localName">{{ region.localName }}</span>
              <span class="text-gray-400 text-xs ml-1 flex-shrink-0">({{ region.regionId }})</span>
            </div>
          </el-checkbox>
        </div>
      </el-checkbox-group>
    </div>
  </div>
</template>

<script setup>
import { reactive, computed, watch } from 'vue'
import { InfoFilled, Refresh } from '@element-plus/icons-vue'
import { getRegions } from '@/plugin/cloud_asset/api/cloudProvider'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  provider: {
    type: String,
    required: true
  },
  accessKey: {
    type: String,
    required: true
  },
  secretKey: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['update:modelValue'])

const regionState = reactive({
  loading: false,
  regions: [],
  selected: [],
  errorMsg: ''
})

// 初始化选中状态
watch(() => props.modelValue, (val) => {
  if (val && regionState.selected.length === 0) {
    regionState.selected = val.split(',')
  }
}, { immediate: true })

const canFetchRegions = computed(() => {
  return props.provider && props.accessKey && props.secretKey
})

const fetchRegions = async () => {
  if (!canFetchRegions.value) return

  regionState.loading = true
  regionState.errorMsg = ''
  regionState.regions = []

  try {
    const res = await getRegions({
      provider: props.provider,
      accessKey: props.accessKey,
      secretKey: props.secretKey
    })

    if (res.code === 0) {
      regionState.regions = res.data || []
    } else {
      regionState.errorMsg = res.msg || '获取可用区失败'
    }
  } catch (error) {
    regionState.errorMsg = '网络请求失败'
    console.error(error)
  } finally {
    regionState.loading = false
  }
}

const handleSelectionChange = (val) => {
  emit('update:modelValue', val.join(','))
}

// 暴露给父组件调用
defineExpose({
  fetchRegions
})
</script>

<style scoped>
/* 深度选择器覆盖 Element Plus 样式以适应 UnoCSS */
:deep(.el-checkbox__label) {
  display: flex;
  align-items: center;
  width: 100%;
  overflow: hidden;
}
</style>
