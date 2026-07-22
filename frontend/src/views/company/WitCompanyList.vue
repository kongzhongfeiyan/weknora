<template>
  <div class="company-list-container">
    <div class="company-list-content">
      <div class="header" style="--wails-draggable: drag">
        <div class="header-title" style="--wails-draggable: drag">
          <div class="title-row" style="--wails-draggable: drag">
            <h2 style="--wails-draggable: drag">{{ $t('witCompany.title') }}</h2>
            <div class="header-actions" style="--wails-draggable: no-drag">
              <t-button theme="primary" size="small" @click="handleCreate">
                <template #icon><t-icon name="add" /></template>
                {{ $t('witCompany.create') }}
              </t-button>
            </div>
          </div>
          <p class="header-subtitle" style="--wails-draggable: drag">{{ $t('witCompany.subtitle') }}</p>
        </div>
      </div>

      <div class="company-list-main">
        <div class="search-bar">
          <t-input
            v-model="searchKeyword"
            :placeholder="$t('witCompany.search')"
            clearable
            @enter="handleSearch"
            @clear="handleSearch"
          >
            <template #prefix-icon><t-icon name="search" /></template>
          </t-input>
        </div>

        <div class="table-container">
          <t-table
            :data="companyList"
            :columns="columns"
            :loading="loading"
            row-key="id"
            :pagination="pagination"
            @page-change="handlePageChange"
            :empty="$t('witCompany.empty')"
          >
            <template #created_at="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
            <template #actions="{ row }">
              <t-space>
                <t-link theme="primary" hover="color" @click="handleEdit(row)">
                  {{ $t('witCompany.edit') }}
                </t-link>
                <t-link theme="danger" hover="color" @click="handleDelete(row)">
                  {{ $t('witCompany.delete') }}
                </t-link>
              </t-space>
            </template>
          </t-table>
        </div>
      </div>
    </div>

    <t-dialog
      v-model:visible="dialogVisible"
      :header="isEdit ? $t('witCompany.edit') : $t('witCompany.create')"
      :confirm-btn="{ loading: submitting }"
      @confirm="handleSubmit"
      @close="handleDialogClose"
      :width="500"
    >
      <t-form ref="formRef" :data="formData" :rules="formRules" layout="vertical">
        <t-form-item :label="$t('witCompany.companyName')" name="company_name">
          <t-input v-model="formData.company_name" :placeholder="$t('witCompany.companyName')" />
        </t-form-item>
        <t-form-item :label="$t('witCompany.companyCode')" name="company_code">
          <t-input v-model="formData.company_code" :placeholder="$t('witCompany.companyCode')" :disabled="isEdit" />
        </t-form-item>
        <t-form-item :label="$t('witCompany.address')" name="address">
          <t-textarea v-model="formData.address" :placeholder="$t('witCompany.address')" :autosize="{ minRows: 2, maxRows: 4 }" />
        </t-form-item>
        <t-form-item :label="$t('witCompany.contactPerson')" name="contact_person">
          <t-input v-model="formData.contact_person" :placeholder="$t('witCompany.contactPerson')" />
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { MessagePlugin, DialogPlugin } from 'tdesign-vue-next'
import type { FormInstanceFunctions, FormRule, PageInfo } from 'tdesign-vue-next'
import { useI18n } from 'vue-i18n'
import {
  listWitCompanies,
  createWitCompany,
  updateWitCompany,
  deleteWitCompany,
  type WitCompany,
  type CreateWitCompanyRequest,
  type UpdateWitCompanyRequest,
} from '@/api/witcompany'

const { t } = useI18n()

const loading = ref(false)
const submitting = ref(false)
const companyList = ref<WitCompany[]>([])
const searchKeyword = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref<FormInstanceFunctions>()

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showPageSize: true,
  showJumper: true,
})

const formData = reactive<CreateWitCompanyRequest>({
  company_name: '',
  company_code: '',
  address: '',
  contact_person: '',
})

const formRules = computed<Record<string, FormRule[]>>(() => ({
  company_name: [{ required: true, message: t('witCompany.companyNameRequired') }],
  company_code: [{ required: true, message: t('witCompany.companyCodeRequired') }],
}))

const columns = computed(() => [
  { colKey: 'company_name', title: t('witCompany.companyName'), width: 200 },
  { colKey: 'company_code', title: t('witCompany.companyCode'), width: 150 },
  { colKey: 'address', title: t('witCompany.address'), ellipsis: true },
  { colKey: 'contact_person', title: t('witCompany.contactPerson'), width: 120 },
  { colKey: 'created_at', title: t('witCompany.createdAt'), width: 180, cell: 'created_at' },
  { colKey: 'actions', title: t('witCompany.actions'), width: 150, cell: 'actions', fixed: 'right' },
])

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString()
}

const fetchList = async () => {
  loading.value = true
  try {
    const res: any = await listWitCompanies(pagination.current, pagination.pageSize, searchKeyword.value || undefined)
    if (res?.success) {
      companyList.value = res.data?.data || []
      pagination.total = res.data?.total || 0
    } else {
      MessagePlugin.error(t('witCompany.loadFailed'))
    }
  } catch {
    MessagePlugin.error(t('witCompany.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  fetchList()
}

const handlePageChange = (pageInfo: PageInfo) => {
  pagination.current = pageInfo.current
  pagination.pageSize = pageInfo.pageSize
  fetchList()
}

const handleCreate = () => {
  isEdit.value = false
  editingId.value = null
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: WitCompany) => {
  isEdit.value = true
  editingId.value = row.id
  formData.company_name = row.company_name
  formData.company_code = row.company_code
  formData.address = row.address
  formData.contact_person = row.contact_person
  dialogVisible.value = true
}

const handleDelete = (row: WitCompany) => {
  const dialog = DialogPlugin.confirm({
    header: t('witCompany.delete'),
    body: t('witCompany.deleteConfirm'),
    confirmBtn: { theme: 'danger' },
    onConfirm: async () => {
      try {
        const res: any = await deleteWitCompany(row.id)
        if (res?.success) {
          MessagePlugin.success(t('witCompany.deleteSuccess'))
          fetchList()
        } else {
          MessagePlugin.error(t('witCompany.deleteFailed'))
        }
      } catch {
        MessagePlugin.error(t('witCompany.deleteFailed'))
      }
      dialog.destroy()
    },
    onClose: () => dialog.destroy(),
  })
}

const handleSubmit = async () => {
  const valid = await formRef.value?.validate()
  if (valid !== true) return

  submitting.value = true
  try {
    let res: any
    if (isEdit.value && editingId.value !== null) {
      const updateData: UpdateWitCompanyRequest = {
        company_name: formData.company_name,
        address: formData.address,
        contact_person: formData.contact_person,
      }
      res = await updateWitCompany(editingId.value, updateData)
      if (res?.success) {
        MessagePlugin.success(t('witCompany.updateSuccess'))
      } else {
        MessagePlugin.error(res?.error?.message || t('witCompany.updateFailed'))
        return
      }
    } else {
      res = await createWitCompany(formData)
      if (res?.success) {
        MessagePlugin.success(t('witCompany.createSuccess'))
      } else {
        MessagePlugin.error(res?.error?.message || t('witCompany.createFailed'))
        return
      }
    }
    dialogVisible.value = false
    fetchList()
  } catch (err: any) {
    MessagePlugin.error(err?.message || (isEdit.value ? t('witCompany.updateFailed') : t('witCompany.createFailed')))
  } finally {
    submitting.value = false
  }
}

const handleDialogClose = () => {
  resetForm()
}

const resetForm = () => {
  formData.company_name = ''
  formData.company_code = ''
  formData.address = ''
  formData.contact_person = ''
  formRef.value?.reset()
}

onMounted(() => {
  fetchList()
})
</script>

<style lang="less" scoped>
.company-list-container {
  display: flex;
  width: 100%;
  height: 100%;
  min-height: 0;
}

.company-list-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 0;
  overflow: hidden;
}

.header {
  flex-shrink: 0;
  padding: 20px 24px 16px;
  border-bottom: 1px solid var(--td-component-stroke);
}

.header-title {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.title-row h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--td-text-color-primary);
}

.header-subtitle {
  margin: 0;
  font-size: 13px;
  color: var(--td-text-color-secondary);
}

.company-list-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
  padding: 16px 24px;
}

.search-bar {
  flex-shrink: 0;
  margin-bottom: 16px;
  max-width: 400px;
}

.table-container {
  flex: 1;
  min-height: 0;
  overflow: auto;
}
</style>
