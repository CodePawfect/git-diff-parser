package parse

import (
	"github.com/codepawfect/git-diff-parser/pkg/model"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	gitDiff := `diff --git a/src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java b/src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java
index 50e23fd0..2b304ea7 100644
--- a/src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java
+++ b/src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java
@@ -10,6 +10,7 @@ import lombok.Data;
 @Data
 @Builder
 public class GoodsReceiptDto {
+  private Long goodsReceiptId;
   private String internalOrderNumber;
   private String OrderNumber;
   private String deliveryNoteNumber;
@@ -26,6 +27,7 @@ public class GoodsReceiptDto {
    */
   public static GoodsReceiptDto toDto(GoodsReceipt goodsReceipt) {
     return GoodsReceiptDto.builder()
+        .goodsReceiptId(goodsReceipt.getId())
         .internalOrderNumber(goodsReceipt.getInternalOrderNumber())
         .OrderNumber(goodsReceipt.getOrderNumber())
         .deliveryNoteNumber(goodsReceipt.getDeliveryNoteNumber())
diff --git a/src/main/java/com/hexagon/adapter/logistic/out/persistence/GoodsReceiptPersistenceAdapter.java b/src/main/java/com/hexagon/adapter/logistic/out/persistence/GoodsReceiptPersistenceAdapter.java
index 13ef59a7..d0c03386 100644
--- a/src/main/java/com/hexagon/adapter/logistic/out/persistence/GoodsReceiptPersistenceAdapter.java
+++ b/src/main/java/com/hexagon/adapter/logistic/out/persistence/GoodsReceiptPersistenceAdapter.java
@@ -45,10 +45,8 @@ public class GoodsReceiptPersistenceAdapter implements GoodsReceiptPersistencePo
   }
 
   @Override
-  public Mono<GoodsReceipt> getByOrderNumber(String OrderNumber) {
-    return goodsReceiptRepository
-        .findByOrderNumber(OrderNumber)
-        .map(GoodsReceiptEntity::mapToDomain);
+  public Mono<GoodsReceipt> getByGoodsReceiptId(Long goodsReceiptId) {
+    return goodsReceiptRepository.findById(goodsReceiptId).map(GoodsReceiptEntity::mapToDomain);
   }`

	result, err := Parse(gitDiff)

	if err != nil {
		t.Fatalf("Parse returned an error: %v", err)
	}

	expected := model.GitDiff{
		FileDiffs: []model.FileDiff{
			{
				OldFilename: "src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java",
				NewFilename: "src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java",
				Hunks: []model.Hunk{
					{
						HunkOperation:    model.ADD,
						OldFileLineStart: 10,
						OldFileLineCount: 6,
						NewFileLineStart: 10,
						NewFileLineCount: 7,
						ChangedLines: []model.ChangedLine{
							{
								Content:    "private Long goodsReceiptId;",
								IsDeletion: false,
							},
						},
					},
					{
						HunkOperation:    model.ADD,
						OldFileLineStart: 26,
						OldFileLineCount: 6,
						NewFileLineStart: 27,
						NewFileLineCount: 7,
						ChangedLines: []model.ChangedLine{
							{
								Content:    ".goodsReceiptId(goodsReceipt.getId())",
								IsDeletion: false,
							},
						},
					},
				},
			},
			{
				OldFilename: "src/main/java/com/hexagon/adapter/logistic/out/persistence/GoodsReceiptPersistenceAdapter.java",
				NewFilename: "src/main/java/com/hexagon/adapter/logistic/out/persistence/GoodsReceiptPersistenceAdapter.java",
				Hunks: []model.Hunk{
					{
						HunkOperation:    model.MODIFY,
						OldFileLineStart: 45,
						OldFileLineCount: 10,
						NewFileLineStart: 45,
						NewFileLineCount: 8,
						ChangedLines: []model.ChangedLine{
							{
								Content:    "public Mono<GoodsReceipt> getByOrderNumber(String OrderNumber) {",
								IsDeletion: true,
							},
							{
								Content:    "return goodsReceiptRepository",
								IsDeletion: true,
							},
							{
								Content:    ".findByOrderNumber(OrderNumber)",
								IsDeletion: true,
							},
							{
								Content:    ".map(GoodsReceiptEntity::mapToDomain);",
								IsDeletion: true,
							},
							{
								Content:    "public Mono<GoodsReceipt> getByGoodsReceiptId(Long goodsReceiptId) {",
								IsDeletion: false,
							},
							{
								Content:    "return goodsReceiptRepository.findById(goodsReceiptId).map(GoodsReceiptEntity::mapToDomain);",
								IsDeletion: false,
							},
						},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("GitDiff is not parsed correctly.\n actual: %v\n expected: %v", result, expected)
	}
}

func TestExtractOldFileName(t *testing.T) {
	input := `diff --git a/src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java b/src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java
	index 50e23fd0..2b304ea7 100644
	--- a/src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java
	+++ b/src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto2.java
	@@ -10,6 +10,7 @@ import lombok.Data;
	@Data
	@Builder
	public class GoodsReceiptDto {
		+  private Long goodsReceiptId;
		private String internalOrderNumber;
		private String OrderNumber;
		private String deliveryNoteNumber;
		@@ -26,6 +27,7 @@ public class GoodsReceiptDto {
		*/
		public static GoodsReceiptDto toDto(GoodsReceipt goodsReceipt) {
		return GoodsReceiptDto.builder()
		+        .goodsReceiptId(goodsReceipt.getId())
		.internalOrderNumber(goodsReceipt.getInternalOrderNumber())
		.OrderNumber(goodsReceipt.getOrderNumber())
		.deliveryNoteNumber(goodsReceipt.getDeliveryNoteNumber())`

	result := extractOldFilename(input)
	expected := "src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java"

	if result != expected {
		t.Errorf("Expected result and expected to be equal, but they are not. result: %s, expected: %s", result, expected)
	}
}

func TestExtractNewFileName(t *testing.T) {
	input := `diff --git a/src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java b/src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java
	index 50e23fd0..2b304ea7 100644
	--- a/src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java
	+++ b/src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto2.java
	@@ -10,6 +10,7 @@ import lombok.Data;
	@Data
	@Builder
	public class GoodsReceiptDto {
		+  private Long goodsReceiptId;
		private String internalOrderNumber;
		private String OrderNumber;
		private String deliveryNoteNumber;
		@@ -26,6 +27,7 @@ public class GoodsReceiptDto {
		*/
		public static GoodsReceiptDto toDto(GoodsReceipt goodsReceipt) {
		return GoodsReceiptDto.builder()
		+        .goodsReceiptId(goodsReceipt.getId())
		.internalOrderNumber(goodsReceipt.getInternalOrderNumber())
		.OrderNumber(goodsReceipt.getOrderNumber())
		.deliveryNoteNumber(goodsReceipt.getDeliveryNoteNumber())`

	result := extractNewFilename(input)
	expected := "src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto2.java"

	if result != expected {
		t.Errorf("Expected result and expected to be equal, but they are not. result: %s, expected: %s", result, expected)
	}
}

func TestExtractChangedLines(t *testing.T) {
	input := `@@ -45,10 +45,8 @@ public class GoodsReceiptPersistenceAdapter implements GoodsReceiptPersistencePo
   }
 
   @Override
-  public Mono<GoodsReceipt> getByOrderNumber(String OrderNumber) {
-    return goodsReceiptRepository
-        .findByOrderNumber(OrderNumber)
-        .map(GoodsReceiptEntity::mapToDomain);
+  public Mono<GoodsReceipt> getByGoodsReceiptId(Long goodsReceiptId) {
+    return goodsReceiptRepository.findById(goodsReceiptId).map(GoodsReceiptEntity::mapToDomain);
   }`

	result := extractChangedLines(input)

	expected := []model.ChangedLine{
		{
			Content:    "public Mono<GoodsReceipt> getByOrderNumber(String OrderNumber) {",
			IsDeletion: true,
		},
		{
			Content:    "return goodsReceiptRepository",
			IsDeletion: true,
		},
		{
			Content:    ".findByOrderNumber(OrderNumber)",
			IsDeletion: true,
		},
		{
			Content:    ".map(GoodsReceiptEntity::mapToDomain);",
			IsDeletion: true,
		},
		{
			Content:    "public Mono<GoodsReceipt> getByGoodsReceiptId(Long goodsReceiptId) {",
			IsDeletion: false,
		},
		{
			Content:    "return goodsReceiptRepository.findById(goodsReceiptId).map(GoodsReceiptEntity::mapToDomain);",
			IsDeletion: false,
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Actual ChangedLines are not equal expected ChangedLines \n actual: %v\n expected: %v", result, expected)
	}
}

func TestExtractHunks(t *testing.T) {
	input := `diff --git a/src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java b/src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java
index 50e23fd0..2b304ea7 100644
--- a/src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java
+++ b/src/main/java/com/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java
@@ -10,6 +10,7 @@ import lombok.Data;
 @Data
 @Builder
 public class GoodsReceiptDto {
+  private Long goodsReceiptId;
   private String internalOrderNumber;
   private String OrderNumber;
   private String deliveryNoteNumber;
@@ -26,6 +27,7 @@ public class GoodsReceiptDto {
    */
   public static GoodsReceiptDto toDto(GoodsReceipt goodsReceipt) {
     return GoodsReceiptDto.builder()
+        .goodsReceiptId(goodsReceipt.getId())
         .internalOrderNumber(goodsReceipt.getInternalOrderNumber())
         .OrderNumber(goodsReceipt.getOrderNumber())
         .deliveryNoteNumber(goodsReceipt.getDeliveryNoteNumber())`

	result, err := extractHunks(input)

	if err != nil {
		t.Fatalf("extract hunks returned an error: %v", err)
	}

	expected := []model.Hunk{
		{
			HunkOperation:    model.ADD,
			OldFileLineStart: 10,
			OldFileLineCount: 6,
			NewFileLineStart: 10,
			NewFileLineCount: 7,
			ChangedLines: []model.ChangedLine{
				{
					Content:    "private Long goodsReceiptId;",
					IsDeletion: false,
				},
			},
		},
		{
			HunkOperation:    model.ADD,
			OldFileLineStart: 26,
			OldFileLineCount: 6,
			NewFileLineStart: 27,
			NewFileLineCount: 7,
			ChangedLines: []model.ChangedLine{
				{
					Content:    ".goodsReceiptId(goodsReceipt.getId())",
					IsDeletion: false,
				},
			},
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Actual Hunks are not equal expected Hunks \n actual: %v\n expected: %v", result, expected)
	}
}
