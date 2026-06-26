export interface Scan {
  id: string
  org_id: string
  type: 'code' | 'webapp'
  status: 'queued' | 'running' | 'completed' | 'failed' | 'cancelled'
  target: string
  created_at: string
  completed_at?: string
}

export interface LLMPentestRun {
  id: string
  org_id: string
  target_endpoint: string
  status: 'queued' | 'running' | 'completed' | 'failed'
  test_modules: string[]
  created_at: string
  completed_at?: string
}

export interface Finding {
  id: string
  org_id: string
  scan_id?: string
  run_id?: string
  title: string
  severity: 'critical' | 'high' | 'medium' | 'low' | 'info'
  description: string
  remediation: string
  created_at: string
}

export interface Report {
  id: string
  org_id: string
  scan_id?: string
  run_id?: string
  format: 'html' | 'md' | 'sarif'
  storage_key: string
  created_at: string
}

export interface Integration {
  id: string
  org_id: string
  provider: 'github'
  created_at: string
}

export interface AlertRule {
  id: string
  org_id: string
  severity_threshold: 'critical' | 'high' | 'medium' | 'low'
  channel: 'email' | 'slack' | 'webhook'
  channel_config: Record<string, unknown>
  active: boolean
  created_at: string
}