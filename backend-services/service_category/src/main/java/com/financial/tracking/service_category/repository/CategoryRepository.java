package com.financial.tracking.service_category.repository;

import java.time.LocalDateTime;
import java.util.UUID;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import com.financial.tracking.service_category.model.CategoryModel;

import jakarta.transaction.Transactional;

@Repository
public interface CategoryRepository extends JpaRepository<CategoryModel, UUID> {
    
    @Modifying
    @Transactional
    @Query(value = """
        INSERT INTO categories_hist (id, name, type, created_at, updated_at)
        SELECT x.id, x.name, x.type::transaction_type, x.created_at, :updated_at
        FROM categories x WHERE x.id = :id
        """, nativeQuery = true)
    int moveToHist(@Param("id") UUID id, @Param("updated_at") LocalDateTime updatedAt);
}

// public interface CategoryHistRepository extends JpaRepository<CategoryModel, UUID> {
    
// }
