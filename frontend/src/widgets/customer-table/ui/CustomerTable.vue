<template>
  <section class="panel">
    <div class="toolbar">
      <input placeholder="Buscar por nome, documento ou e-mail" />
      <button>Exportar</button>
    </div>
    <table>
      <thead>
        <tr>
          <th>Nome</th>
          <th>Documento</th>
          <th>E-mail</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>Padaria Central Ltda</td>
          <td>12.345.678/0001-90</td>
          <td>contato@padariacentral.com.br</td>
          <td><span class="badge">Ativo</span></td>
        </tr>
      </tbody>
    </table>
  </section>
</template>

<style scoped>
.panel {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 10px;
  overflow: hidden;
}

.toolbar {
  display: flex;
  gap: 10px;
  padding: 12px;
}

input {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  flex: 1;
  padding: 8px 12px;
}

table {
  border-collapse: collapse;
  width: 100%;
}

th,
td {
  border-top: 1px solid var(--color-border);
  font-size: 13px;
  padding: 12px 16px;
  text-align: left;
}

.badge {
  background: #dcfce7;
  border-radius: 999px;
  color: #15803d;
  padding: 3px 10px;
}
</style>
