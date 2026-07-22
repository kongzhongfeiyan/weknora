import { get, post, put, del } from '@/utils/request'

export interface WitCompany {
  id: number
  company_name: string
  company_code: string
  address: string
  contact_person: string
  created_at: string
  created_by: string
}

export interface CreateWitCompanyRequest {
  company_name: string
  company_code: string
  address?: string
  contact_person?: string
  created_by?: string
}

export interface UpdateWitCompanyRequest {
  company_name?: string
  company_code?: string
  address?: string
  contact_person?: string
}

export const listWitCompanies = (page: number, pageSize: number, keyword?: string) =>
  get(`/api/v1/witcompanies?page=${page}&page_size=${pageSize}${keyword ? `&keyword=${encodeURIComponent(keyword)}` : ''}`)

export const getWitCompanyById = (id: number) =>
  get(`/api/v1/witcompanies/${id}`)

export const createWitCompany = (data: CreateWitCompanyRequest) =>
  post('/api/v1/witcompanies', data)

export const updateWitCompany = (id: number, data: UpdateWitCompanyRequest) =>
  put(`/api/v1/witcompanies/${id}`, data)

export const deleteWitCompany = (id: number) =>
  del(`/api/v1/witcompanies/${id}`)
