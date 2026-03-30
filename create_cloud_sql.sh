#!/bin/bash

# Configuration Variables
PROJECT_ID="dev-design-491813"        # Project ID จาก gcloud projects list
INSTANCE_NAME="ticket-booking-db"     # ชื่อ Database Instance ที่จะสร้าง
REGION="asia-southeast1"              # โซนสิงคโปร์ (ใกล้ไทยที่สุด)
TIER="db-f1-micro"                    # สเปคเครื่อง: f1-micro ประหยัดสุด เหมาะสำหรับ Dev/Test
DATABASE_VERSION="POSTGRES_15"        # เวอร์ชัน PostgreSQL
DB_USER="postgres"
DB_PASS="StrongSecret123!"            # รหัสผ่านสำหรับ Database (เปลี่ยนก่อนรัน!)

echo "====================================================="
echo "🗄️  Creating Cloud SQL Instance: $INSTANCE_NAME"
echo "Project: $PROJECT_ID, Region: $REGION, Tier: $TIER"
echo "====================================================="

# 1. Set Active Project
echo "👉 Setting gcloud project to $PROJECT_ID..."
gcloud config set project "$PROJECT_ID"

# 2. Enable Cloud SQL API (ถ้ายังไม่ได้เปิด)
echo "👉 Enabling Cloud SQL Admin API..."
gcloud services enable sqladmin.googleapis.com

# 3. Create Instance
echo "👉 Creating Cloud SQL $DATABASE_VERSION instance (This will take a few minutes)..."
gcloud sql instances create "$INSTANCE_NAME" \
    --database-version="$DATABASE_VERSION" \
    --tier="$TIER" \
    --region="$REGION" \
    --project="$PROJECT_ID" \
    --root-password="$DB_PASS"

# 4. (Optional) Create Logical Database Name inside the instance
echo "👉 Creating main database inside the instance..."
gcloud sql databases create tickets_db --instance="$INSTANCE_NAME"

# 5. Get Connection Name
CONNECTION_NAME=$(gcloud sql instances describe "$INSTANCE_NAME" --format="value(connectionName)")

echo "====================================================="
echo "✅ Cloud SQL created successfully!"
echo "Connection Name: $CONNECTION_NAME"
echo "Username: $DB_USER"
echo "Password: $DB_PASS"
echo "Connection String format for local dev: "
echo "host=localhost user=$DB_USER password=$DB_PASS dbname=tickets_db port=5432 sslmode=disable"
echo "====================================================="
echo "📌 ถ้าจะต่อผ่าน Local อย่าลืมติดตั้ง Cloud SQL Auth Proxy แล้วรัน:"
echo "./cloud-sql-proxy $CONNECTION_NAME"
echo "====================================================="
