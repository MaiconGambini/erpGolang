export interface ApiResponse<T> {
  data: T
}

export interface ApiErrorBody {
  error: {
    code: string
    message: string
    details?: Record<string, string>
  }
}

export interface Paginated<T> {
  data: T[]
  pagination: {
    limit: number
    offset: number
    total: number
  }
}
