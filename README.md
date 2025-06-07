# README

## About

# Realtor Installment Assistant

A cross-platform tool for managing property sales in installments, tracking payments, buyers, commissions, and generating key reports — designed to run locally (SQLite + Go + Gin) today and easily migrate to serverless/cloud or a Fyne desktop UI tomorrow.

---

## 🔍 Project Overview

**Realtor Installment Assistant** helps real-estate sellers:

- **Record properties** along with location-based pricing  
- **Create installment plans** (amount due, schedule) per buyer & property  
- **Track individual payments** and outstanding balances  
- **Manage buyers, lettings, sales, and introductions**  
- **Compute & record commissions** (sales, lettings, introductions) to any beneficiary  
- **Dynamically enforce permissions** via roles & permissions tables  
- **Expose reporting endpoints** returning JSON summaries for:
  - Total commissions by beneficiary  
  - Outstanding installments by plan  
  - Monthly sales volume  
  - Active lettings rent roll  
  - Top properties by payment volume  

Future directions include a **Fyne desktop front-end**, optional cloud deployment, automated migrations, CSV/XLSX exports, and tenant-based licensing.

---

## 🚀 Tech Stack

- **Backend**: Go 1.20+, [Gin](https://github.com/gin-gonic/gin) HTTP framework  
- **Persistence**: SQLite (one `.db` per domain), in-code migrations  
- **Auth & RBAC**: JWT (`github.com/golang-jwt/jwt/v4`), dynamic roles & permissions  
- **Repo pattern**: One Go repository per entity (Property, Buyer, Commission, etc.)  
- **Reporting**: Direct SQL → Go structs; no database views  
- **Future UI**: Fyne (desktop) or React/Electron (web)  

---

## ⚙️ Getting Started

1. **Clone & build**  
   ```bash
   git clone https://github.com/your-org/realtor-installment.git
   cd realtor-installment/api
   go mod tidy
   go run main.go
````

2. **Environment**

   ```bash
   export APP_JWT_SECRET="a-long-random-secret-string"
   ```
3. **API Endpoints**

   * `POST /login` → `{ username, password }` → `{ token }`
   * `POST /register` → admin only
   * CRUD for `/properties`, `/buyers`, `/plans`, `/installments`, `/payments`, `/commissions`, `/sales`, `/lettings`, `/introductions`, `/users`
   * Reporting:

     * `GET /reports/commissions/beneficiary`
     * `GET /reports/installments/outstanding`
     * `GET /reports/sales/monthly`
     * `GET /reports/lettings/rentroll`
     * `GET /reports/properties/top-payments`

   All protected by JWT + permission checks.

---

## 🛠️ Development

* **Migrations** live in `dbmigrations/migrations.go`; they run automatically on startup.
* **Repo implementations** in `api/repos/sqlite_*.go`.
* **Services** in `api/services/*.go`; each encapsulates business logic.
* **Handlers** in `api/handlers/*.go`; each ties service → HTTP.
* **Middleware** (`AuthMiddleware`, `RequirePermission`) in `main.go`.

To add a new report, simply:

1. Write a **SQL query** in the relevant `sqlite_*.go` repo.
2. Add a **service** method in `report_service.go`.
3. Add a **handler** method in `report.go`.
4. Wire up a new **route** in `main.go`.

---

## 🔮 Roadmap & Future Work

* **Fyne Desktop App** with local SQLite and optional license file
* **Cloud Deployment** on serverless (AWS Lambda / Azure Functions)
* **CSV / Excel exports** and scheduled report generation
* **Refresh tokens**, **password reset**, **MFA** (if needed)
* **Multi-tenant licensing** with expiry and feature tiers

---

## 📄 License

This project is released under the MIT License. See [LICENSE](LICENSE) for details.

---

> Feel free to ⭐ star, fork & contribute!

```
