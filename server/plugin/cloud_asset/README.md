# 云资产管理 (Cloud Asset Management) 开发文档

## 1. 系统架构重构说明

### 1.1 菜单重构
- **重构内容**：将原属于“可观测可视化”子菜单的“项目管理”模块提升为一级菜单，并更名为“前端可观测可视化项目管理”。
- **目的**：提高项目管理模块的可见性，使其成为系统中通用的资源管理单元，便于被“云资产管理”等其他业务模块引用。
- **权限控制**：该菜单项在 `sys_base_menus` 中 `parent_id` 为 `0`，可通过角色管理进行独立的权限分配。

### 1.2 数据关联机制
- **核心逻辑**：云厂商 (CloudProvider) 不再直接存储厂商名称，而是通过 `ProjectID` 关联“项目管理”模块中的项目。
- **数据一致性**：采用数据库 ID 关联（Foreign Key 逻辑）。后端在查询云厂商列表时，通过 `Preload("Project")` 实时获取最新的项目名称。
- **实时同步**：当“项目管理”中的项目名称修改时，由于云厂商页面是基于 ID 动态加载数据的，厂商名称会实现自动同步更新。

## 2. 后端接口规范

### 2.1 云厂商管理接口
所有接口前缀为 `/cloudProvider/`。

| 接口名 | 方法 | 功能描述 | 关联项说明 |
| :--- | :--- | :--- | :--- |
| `createCloudProvider` | POST | 创建云厂商 | 需传 `projectId` |
| `updateCloudProvider` | PUT | 更新云厂商 | 支持修改关联项目 |
| `deleteCloudProvider` | DELETE | 删除云厂商 | 物理/逻辑删除 |
| `getCloudProviderList` | POST | 分页获取列表 | 返回数据包含 `project` 对象 |

### 2.2 数据共享机制
- **共享接口**：通过引用 `cert_manager` 插件提供的 `getCertCategoryList` 接口获取项目候选列表。
- **规范**：任何需要引用项目的模块，应统一使用 `ProjectID` 进行关联，并确保在 Service 层使用预加载机制。

## 3. 前端组件使用说明

### 3.1 云厂商页面 (`cloudProvider.vue`)
- **功能特性**：
    - **动态关联**：新增/编辑时，厂商名称为禁用状态，用户需从“关联项目”下拉框中选择项目。
    - **自动同步**：选择项目后，前端通过 `handleProjectChange` 自动填充厂商名称。
    - **脱敏显示**：AK/SK 在列表中经过掩码处理，仅在编辑模式下可点击查看。

### 3.2 引用示例
```javascript
import { getCertCategoryList } from '@/plugin/cert_manager/api/certCategory'

// 获取项目列表用于下拉选择
const getProjectList = async () => {
  const res = await getCertCategoryList({ page: 1, pageSize: 999 })
  if (res.code === 0) {
    projectList.value = res.data.list
  }
}
```

## 4. 测试验证
- **单元测试**：见 `server/plugin/cloud_asset/service/cloud_provider_test.go`。
- **验证内容**：
    1. 验证云厂商创建及项目关联功能。
    2. 验证项目更名后，云厂商列表数据的级联同步展示。
    3. 验证列表分页及筛选功能。
