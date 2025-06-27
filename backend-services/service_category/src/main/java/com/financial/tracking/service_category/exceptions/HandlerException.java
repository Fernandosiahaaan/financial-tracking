package com.financial.tracking.service_category.exceptions;

import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.*;

import com.financial.tracking.service_category.dto.BaseResponse;

import java.util.stream.Collectors;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

@RestControllerAdvice
public class HandlerException {
      @ExceptionHandler(NotFoundException.class)
      public ResponseEntity<BaseResponse<Object>> handleNotFound(NotFoundException ex){
        return ResponseEntity.status(HttpStatus.NOT_FOUND).body(BaseResponse.setResponse(false, ex.getMessage(), "[NotFoundException] " + ex.getMessage(), null));
      } 

      @ExceptionHandler(MethodArgumentNotValidException.class)
      public ResponseEntity<BaseResponse<Void>> handleValidationException(MethodArgumentNotValidException ex) {
          String messageErr = ex.getBindingResult().getFieldErrors()
              .stream()
              .map(e -> e.getField() + ": " + e.getDefaultMessage())
              .collect(Collectors.joining("; "));
          return ResponseEntity.badRequest().body(BaseResponse.setResponse(false, "kesalahan pada sistem, mohon menghubungi petugas [E002]", "[ValidationException] " + messageErr, null));
      }

      @ExceptionHandler(BadRequestException.class)
      public ResponseEntity<BaseResponse<Object>> handleBadRequest(BadRequestException ex) {
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(BaseResponse.setResponse(false, ex.getMessage(), "[BadRequestException] "+ ex.getMessage(), null));
      }

      @ExceptionHandler(RuntimeException.class)
      public ResponseEntity<BaseResponse<Object>> handleRunTime(RuntimeException ex){
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(BaseResponse.setResponse(false, "kesalahan pada sistem, mohon menghubungi petugas [E004]", "[INTERNAL_SERVER_ERROR] : " + ex.getMessage(), null));
      }

}
