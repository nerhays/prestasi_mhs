// database/mongo-init.js
// Run mongosh < database/mongo-init.js

db = db.getSiblingDB("prestasi_db");

db.createCollection("achievements", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["studentId", "achievementType", "title", "description", "details", "points", "createdAt", "updatedAt"],
      properties: {
        studentId: { bsonType: "string", description: "UUID mahasiswa (refer ke PostgreSQL students.id)" },
        achievementType: {
          bsonType: "string",
          enum: ["academic", "competition", "organization", "publication", "certification", "other"],
          description: "tipe prestasi",
        },
        title: { bsonType: "string" },
        description: { bsonType: "string" },
        details: {
          bsonType: "object",
          description: "field dinamis sesuai tipe prestasi",
          properties: {
            competitionName: { bsonType: "string" },
            competitionLevel: { bsonType: "string" },
            rank: { bsonType: ["int", "long", "double"] },
            medalType: { bsonType: "string" },

            publicationType: { bsonType: "string" },
            publicationTitle: { bsonType: "string" },
            authors: { bsonType: "array", items: { bsonType: "string" } },
            publisher: { bsonType: "string" },
            issn: { bsonType: "string" },

            organizationName: { bsonType: "string" },
            position: { bsonType: "string" },
            period: {
              bsonType: "object",
              properties: {
                start: { bsonType: "date" },
                end: { bsonType: "date" },
              },
            },

            certificationName: { bsonType: "string" },
            issuedBy: { bsonType: "string" },
            certificationNumber: { bsonType: "string" },
            validUntil: { bsonType: "date" },

            eventDate: { bsonType: "date" },
            location: { bsonType: "string" },
            organizer: { bsonType: "string" },
            score: { bsonType: ["int", "long", "double"] },
            customFields: { bsonType: "object" },
          },
        },
        attachments: {
          bsonType: "array",
          items: {
            bsonType: "object",
            required: ["fileName", "fileUrl", "fileType", "uploadedAt"],
            properties: {
              fileName: { bsonType: "string" },
              fileUrl: { bsonType: "string" },
              fileType: { bsonType: "string" },
              uploadedAt: { bsonType: "date" },
            },
          },
        },
        tags: { bsonType: "array", items: { bsonType: "string" } },
        points: { bsonType: ["int", "long", "double"] },
        createdAt: { bsonType: "date" },
        updatedAt: { bsonType: "date" },
        isDeleted: { bsonType: "bool" },
        deletedAt: { bsonType: "date" },
      },
    },
  },
  validationLevel: "moderate",
});

db.achievements.createIndex({ studentId: 1 });
db.achievements.createIndex({ achievementType: 1 });
db.achievements.createIndex({ "details.competitionLevel": 1 });
db.achievements.createIndex({ createdAt: -1 });
