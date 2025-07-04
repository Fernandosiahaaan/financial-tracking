package com.financial.tracking.service_category.services;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.financial.tracking.service_category.dto.BaseResponse;
import com.financial.tracking.service_category.dto.CategoryRequest;
import com.financial.tracking.service_category.exceptions.*;
import com.financial.tracking.service_category.model.CategoryModel;
import com.financial.tracking.service_category.repository.CategoryRepository;

@Service
public class CategoryServiceImpl implements CategoryService {
    @Autowired
    private CategoryRepository categoryRepository;
    
    @Override
    public BaseResponse<CategoryModel> create(CategoryRequest item){
        if (categoryRepository.findByName(item.getName()).isPresent()) {
            throw new BadRequestException("Category name '" + item.getName() + "' sudah tersedia");
        }

        CategoryModel data = new CategoryModel();
        data.setId(UUID.randomUUID());
        data.setName(item.getName());
        data.setType(item.getType());
        data.setCreatedAt(LocalDateTime.now());
        CategoryModel dataSave =  categoryRepository.save(data);
        return BaseResponse.setResponse(true, "Success", null, dataSave);
    }

    public BaseResponse<List<CategoryModel>> findAll() {
        List<CategoryModel> datas = categoryRepository.findAll();
        return BaseResponse.setResponse(true, "Success", null, datas);
    }

    public BaseResponse<CategoryModel> findById(UUID id) {
        CategoryModel data = categoryRepository.findById(id).orElseThrow(() -> new NotFoundException("Category dengan id '"+ id + "' tidak ditemukan"));
        return BaseResponse.setResponse(true, "Success", null, data);
    }

    public BaseResponse<CategoryModel> update(CategoryRequest item) {
        CategoryModel existing = categoryRepository.findById(item.getId()).orElseThrow(() -> new NotFoundException("Category dengan id '"+ item.getId() + "' tidak ditemukan"));

        existing.setName(item.getName());
        existing.setType(item.getType());
        existing.setUpdatedAt(LocalDateTime.now());
        CategoryModel data = categoryRepository.save(existing);
        return BaseResponse.setResponse(true, "Success", null, data);
    }

    public BaseResponse<Void> delete(UUID id){
        CategoryModel data = categoryRepository.findById(id).orElseThrow(() -> new NotFoundException("Category dengan id '"+ id + "' tidak ditemukan"));
        System.out.println("Memindahkan ID: " + data.getId());
        int result = categoryRepository.moveToHist(data.getId(), LocalDateTime.now());
        System.out.println("Result dari moveToHist: " + result);
        if (result ==  0) {
            throw new NotFoundException("category with id `" + id.toString() + "' gagal dihapus");
        }
        categoryRepository.deleteById(id);
        return BaseResponse.setResponse(true, "Success delete id '" + id.toString()+ "' ", null, null);
    }
}
