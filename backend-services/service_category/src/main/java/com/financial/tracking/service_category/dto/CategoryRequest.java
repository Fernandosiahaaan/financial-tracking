package com.financial.tracking.service_category.dto;

import java.util.UUID;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.financial.tracking.service_category.model.transactionType;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.*;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder

public class CategoryRequest {
    @JsonProperty("id")
    private UUID id;

    @JsonProperty("name")
    @NotBlank(message = "Name is required")
    private String name;

    @JsonProperty("type")
    @NotNull(message = "Type is required")
    private transactionType type;
}
