package com.tracking.financial.service_transaction.clients;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.tracking.financial.service_transaction.dto.BaseResponse;

@Component
public class ClientCategory {

    @Value("${service-category.url}")
    private String categoryBaseUrl;

    private final HttpClient httpClient = HttpClient.newHttpClient();
    private final ObjectMapper mapper = new ObjectMapper();

    public BaseResponse getCategoryById(String categoryId) {
        String urlString = categoryBaseUrl + "/ctgr/" + categoryId;

        try {
            HttpRequest request = HttpRequest.newBuilder()
                    .uri(URI.create(urlString))
                    .header("Accept", "application/json")
                    .GET()
                    .build();

            HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
            // Parse JSON
            JsonNode root = mapper.readTree(response.body());

            boolean success = root.path("success").asBoolean(false);
            String message = root.path("message").asText(null);
            String messageError = root.path("messageError").asText(null);
            JsonNode dataNode = root.path("data");

            Object data = null;
            if (dataNode.isObject()) {
                data = mapper.convertValue(dataNode, Object.class);
            }

            return BaseResponse.setResponse(success, message, messageError, data);

        } catch (InterruptedException e) {
            Thread.currentThread().interrupt(); // only for InterruptedException
            return BaseResponse.setResponse(false, "Request interrupted", e.getMessage(), null);
        } catch (IOException e) {
            return BaseResponse.setResponse(false, "Request failed to get category by ID", e.getMessage(), null);
        }
    }

    public Boolean isExistCategoryById(String categoryId) {
        BaseResponse response = getCategoryById(categoryId);
        return response.isSuccess();
    }
}