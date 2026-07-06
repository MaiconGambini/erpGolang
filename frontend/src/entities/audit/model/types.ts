export interface AuditLog {
  id: string
  actorUserId?: string
  action: string
  entityType: string
  entityId: string
  createdAt: string
}
