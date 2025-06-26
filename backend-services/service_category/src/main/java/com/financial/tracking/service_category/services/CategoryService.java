package com.financial.tracking.service_category.services;

import java.util.List;
import java.util.UUID;

import com.financial.tracking.service_category.dto.BaseResponse;
import com.financial.tracking.service_category.dto.CategoryRequest;
import com.financial.tracking.service_category.model.CategoryModel;

public interface CategoryService {
    BaseResponse<CategoryModel> create(CategoryRequest item);
    BaseResponse<List<CategoryModel>> findAll();
    BaseResponse<CategoryModel> findById(UUID id);
    BaseResponse<CategoryModel> update(CategoryRequest item);
    BaseResponse<Void> delete(UUID id);
}
