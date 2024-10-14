# git diff parser
A lightweight Go library for parsing Git diffs into structured, easy-to-process data, enabling developers to programmatically analyze and manipulate file changes with minimal effort.

## Example
Demonstration how a `GitDiff` struct might look after parsing a git diff.

### Input: Git Diff

```
diff --git a/src/main/java/com/eon/smexnet/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java b/src/main/java/com/eon/smexnet/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java
index 50e23fd0..2b304ea7 100644
--- a/src/main/java/com/eon/smexnet/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java
+++ b/src/main/java/com/eon/smexnet/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java
@@ -10,6 +10,7 @@ import lombok.Data;
 @Data
 @Builder
 public class GoodsReceiptDto {
+  private Long goodsReceiptId;
   private String internalOrderNumber;
   private String eonOrderNumber;
   private String deliveryNoteNumber;
@@ -26,6 +27,7 @@ public class GoodsReceiptDto {
    */
   public static GoodsReceiptDto toDto(GoodsReceipt goodsReceipt) {
     return GoodsReceiptDto.builder()
+        .goodsReceiptId(goodsReceipt.getId())
         .internalOrderNumber(goodsReceipt.getInternalOrderNumber())
         .eonOrderNumber(goodsReceipt.getEonOrderNumber())
         .deliveryNoteNumber(goodsReceipt.getDeliveryNoteNumber())
diff --git a/src/main/java/com/eon/smexnet/hexagon/adapter/logistic/out/persistence/GoodsReceiptPersistenceAdapter.java b/src/main/java/com/eon/smexnet/hexagon/adapter/logistic/out/persistence/GoodsReceiptPersistenceAdapter.java
index 13ef59a7..d0c03386 100644
--- a/src/main/java/com/eon/smexnet/hexagon/adapter/logistic/out/persistence/GoodsReceiptPersistenceAdapter.java
+++ b/src/main/java/com/eon/smexnet/hexagon/adapter/logistic/out/persistence/GoodsReceiptPersistenceAdapter.java
@@ -45,10 +45,8 @@ public class GoodsReceiptPersistenceAdapter implements GoodsReceiptPersistencePo
   }
 
   @Override
-  public Mono<GoodsReceipt> getByEonOrderNumber(String eonOrderNumber) {
-    return goodsReceiptRepository
-        .findByEonOrderNumber(eonOrderNumber)
-        .map(GoodsReceiptEntity::mapToDomain);
+  public Mono<GoodsReceipt> getByGoodsReceiptId(Long goodsReceiptId) {
+    return goodsReceiptRepository.findById(goodsReceiptId).map(GoodsReceiptEntity::mapToDomain);
   }
```

### Output: Parsed GitDiff Struct

```go
model.GitDiff{
    FileDiffs: []model.FileDiff{
        {
            OldFilename: "src/main/java/com/eon/smexnet/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java",
            NewFilename: "src/main/java/com/eon/smexnet/hexagon/adapter/logistic/in/rest/model/GoodsReceiptDto.java",
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
            OldFilename: "src/main/java/com/eon/smexnet/hexagon/adapter/logistic/out/persistence/GoodsReceiptPersistenceAdapter.java",
            NewFilename: "src/main/java/com/eon/smexnet/hexagon/adapter/logistic/out/persistence/GoodsReceiptPersistenceAdapter.java",
            Hunks: []model.Hunk{
                {
                    HunkOperation:    model.MODIFY,
                    OldFileLineStart: 45,
                    OldFileLineCount: 10,
                    NewFileLineStart: 45,
                    NewFileLineCount: 8,
                    ChangedLines: []model.ChangedLine{
                        {
                            Content:    "public Mono<GoodsReceipt> getByEonOrderNumber(String eonOrderNumber) {",
                            IsDeletion: true,
                        },
                        {
                            Content:    "return goodsReceiptRepository",
                            IsDeletion: true,
                        },
                        {
                            Content:    ".findByEonOrderNumber(eonOrderNumber)",
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
```

## Usage

```bash
go get github.com/codepawfect/git-diff-parser
```

- Input: You pass a Git diff as a string to the Parse function.
- Parsing: The Parse function processes the diff and returns a GitDiff struct, which contains detailed information about file changes, including hunks and lines that were added or deleted.

## License
This project is licensed under the terms of the [MIT License](./LICENSE).