import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/pages/Dashboard.vue')
  },
  {
    path: '/sign-in',
    name: 'sign-in',
    component: () => import('@/views/sign-in.vue')
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
  },
  {
    path: '/docs',
    name: 'docs',
    component: () => import('@/pages/Docs.vue')
  },
  {
    path: '/sign-up',
    name: 'sign-up',
    component: () => import('@/views/sign-up.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router