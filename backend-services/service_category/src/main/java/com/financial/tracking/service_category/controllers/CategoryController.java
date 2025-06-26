package com.financial.tracking.service_category.controllers;

import java.util.UUID;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.financial.tracking.service_category.dto.BaseResponse;
import com.financial.tracking.service_category.dto.CategoryRequest;
import com.financial.tracking.service_category.model.CategoryModel;
import com.financial.tracking.service_category.services.CategoryService;

import jakarta.validation.Valid;

import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.PutMapping;




@RestController
@RequestMapping("/api/categories")
public class CategoryController {
    
    @Autowired
    private CategoryService categoryService;

    @PostMapping
    public ResponseEntity<BaseResponse<?>> create(@Valid @RequestBody CategoryRequest request) {
        return ResponseEntity.status(HttpStatus.OK).body(categoryService.create(request));
    }

    @GetMapping
    public ResponseEntity<BaseResponse<?>> findAll() {
        return ResponseEntity.status(HttpStatus.OK).body(categoryService.findAll());
    }

    @GetMapping("/{id}")
    public ResponseEntity<BaseResponse<?>> findById(@PathVariable UUID id) {
        return ResponseEntity.status(HttpStatus.OK).body(categoryService.findById(id));
    }

    @PutMapping
    public ResponseEntity<BaseResponse<?>> update(@Valid @RequestBody CategoryRequest request) {
        return ResponseEntity.status(HttpStatus.OK).body(categoryService.update(request));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<BaseResponse<?>> delete(@PathVariable UUID id) {
        return ResponseEntity.status(HttpStatus.OK).body(categoryService.delete(id));
    }
    
}
