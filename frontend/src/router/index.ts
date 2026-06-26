import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'dashboard',
    component: () => import('@/pages/Dashboard.vue')
  },
  {
    path: '/scans',
    name: 'scans',
    component: () => import('@/pages/Scans.vue')
  },
  {
    path: '/scans/:id',
    name: 'scan-detail',
    component: () => import('@/pages/ScanDetail.vue')
  },
  {
    path: '/llmpentest',
    name: 'llmpentest',
    component: () => import('@/pages/LLMPentest.vue')
  },
  {
    path: '/llmpentest/:id',
    name: 'llmpentest-detail',
    component: () => import('@/pages/LLMPentestDetail.vue')
  },
  {
    path: '/reports',
    name: 'reports',
    component: () => import('@/pages/Reports.vue')
  },
  {
    path: '/reports/:id',
    name: 'report-detail',
    component: () => import('@/pages/ReportDetail.vue')
  },
  {
    path: '/integrations',
    name: 'integrations',
    component: () => import('@/pages/Integrations.vue')
  },
  {
    path: '/alerts',
    name: 'alerts',
    component: () => import('@/pages/Alerts.vue')
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('@/pages/Settings.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router