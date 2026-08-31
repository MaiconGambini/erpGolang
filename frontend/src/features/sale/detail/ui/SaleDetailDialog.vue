<template>
  <AppDialog :visible="visible" title-id="detail-title" size="lg" @close="emit('close')">
      <h2 id="detail-title">Detalhes da venda</h2>
      <p v-if="isLoading" class="state">Carregando...</p>
      <p v-else-if="isError" class="state error">Erro ao carregar venda</p>
      <template v-else-if="sale">
        <dl class="meta">
          <div>
            <dt>Cliente</dt>
            <dd>{{ sale.customerName }}</dd>
          </div>
          <div>
            <dt>Status</dt>
            <dd>
              <span class="badge" :class="sale.status">{{ statusLabel(sale.status) }}</span>
            </dd>
          </div>
          <div>
            <dt>Total</dt>
            <dd class="total">R$ {{ formatPrice(sale.total) }}</dd>
          </div>
          <div>
            <dt>Data</dt>
            <dd>{{ formatDate(sale.createdAt) }}</dd>
          </div>
        </dl>

        <section v-if="sale.notes" class="notes">
          <h3>Observações</h3>
          <p>{{ sale.notes }}</p>
        </section>

        <section class="items-section">
          <h3>Itens</h3>
          <p v-if="!sale.items?.length" class="state">Nenhum item</p>
          <table v-else class="items-table">
            <thead>
              <tr>
                <th>Produto</th>
                <th>SKU</th>
                <th>Qtd</th>
                <th>Preço unit.</th>
                <th>Subtotal</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in sale.items" :key="item.id">
                <td>{{ item.productName }}</td>
                <td>{{ item.productSku }}</td>
                <td>{{ item.quantity }}</td>
                <td>R$ {{ formatPrice(item.unitPrice) }}</td>
                <td>R$ {{ formatPrice(item.lineTotal) }}</td>
              </tr>
            </tbody>
          </table>
        </section>

        <div class="actions">
          <button type="button" class="secondary" @click="emit('close')">Fechar</button>
          <button type="button" class="secondary" :disabled="downloadingPdf" @click="onDownloadPdf">
            Baixar PDF
          </button>
          <AppButton v-if="sale.status === 'draft' && canWriteUser" @click="emit('edit', sale.id)">
            Editar
          </AppButton>
        </div>
      </template>
  </AppDialog>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { getSale } from '@/entities/sale/api/sale.api'
import { downloadSalePdf } from '@/entities/reports/api/reports.api'
import type { SaleStatus } from '@/entities/sale/model/types'
import { useSessionStore } from '@/entities/session/model/session.store'
import AppDialog from '@/shared/ui/AppDialog.vue'
import AppButton from '@/shared/ui/AppButton.vue'
import { downloadBlob } from '@/shared/lib/download'
import { canWrite } from '@/shared/lib/roles'

const props = defineProps<{ visible: boolean; saleId: string | null }>()
const emit = defineEmits<{ close: []; edit: [id: string] }>()

const session = useSessionStore()
const canWriteUser = computed(() => canWrite(session.user?.role))
const downloadingPdf = ref(false)

const { data: sale, isLoading, isError } = useQuery({
  queryKey: computed(() => ['sale', props.saleId]),
  queryFn: () => getSale(props.saleId!),
  enabled: computed(() => props.visible && !!props.saleId),
})

function formatPrice(value: string) {
  const n = Number(value)
  if (Number.isNaN(n)) return value
  return n.toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

function formatDate(value: string) {
  return new Date(value).toLocaleDateString('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function statusLabel(status: SaleStatus) {
  const map: Record<SaleStatus, string> = {
    draft: 'Rascunho',
    confirmed: 'Confirmada',
    cancelled: 'Cancelada',
  }
  return map[status]
}

async function onDownloadPdf() {
  if (!props.saleId) return
  downloadingPdf.value = true
  try {
    const blob = await downloadSalePdf(props.saleId)
    downloadBlob(`venda-${props.saleId}.pdf`, blob)
  } finally {
    downloadingPdf.value = false
  }
}
</script>

<style scoped>

h2 {
  font-size: 18px;
  font-weight: 650;
  margin: 0 0 20px;
}

h3 {
  font-size: 14px;
  margin: 0 0 10px;
}

.state {
  color: var(--color-text-muted);
  font-size: 14px;
}

.state.error {
  color: var(--color-danger);
}

.meta {
  display: grid;
  gap: 12px;
  grid-template-columns: 1fr 1fr;
  margin: 0 0 20px;
}

.meta div {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

dt {
  color: var(--color-text-secondary);
  font-size: 12px;
}

dd {
  font-size: 14px;
  margin: 0;
}

.total {
  font-weight: 600;
}

.badge {
  border-radius: 999px;
  display: inline-block;
  font-size: 12px;
  font-weight: 650;
  line-height: 1;
  padding: 6px 9px;
  width: fit-content;
}

.badge.draft {
  background: var(--color-warning-soft);
  color: var(--color-warning-strong);
}

.badge.confirmed {
  background: var(--color-success-soft);
  color: var(--color-success-strong);
}

.badge.cancelled {
  background: var(--color-surface-muted);
  color: var(--color-text-muted);
}

.notes {
  margin-bottom: 20px;
}

.notes p {
  color: var(--color-text-secondary);
  font-size: 14px;
  margin: 0;
  white-space: pre-wrap;
}

.items-section {
  margin-bottom: 20px;
}

.items-table {
  border-collapse: collapse;
  width: 100%;
}

.items-table th,
.items-table td {
  border-top: 1px solid var(--color-border);
  font-size: 13px;
  padding: 8px 10px;
  text-align: left;
}

.items-table th {
  color: var(--color-text-secondary);
  font-weight: 500;
}

.actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 8px;
}

.secondary {
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  font: inherit;
  padding: 8px 14px;
}

.secondary:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}
</style>
