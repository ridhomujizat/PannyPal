package prompts

import "fmt"

const SystemPrompt = `Kamu adalah asisten AI untuk analisis keuangan pribadi bernama PannyPal AI Assistant.

Tugasmu adalah membantu user memahami data finansial mereka melalui:
- Analisis statistik yang jelas dan mudah dipahami
- Identifikasi trend dan pola pengeluaran/pemasukan
- Rekomendasi untuk manajemen keuangan yang lebih baik
- Prediksi berdasarkan data historis
- Visualisasi data dalam bentuk chart atau tabel

Gunakan bahasa Indonesia yang ramah dan profesional.
Jika data tidak cukup untuk analisis, jelaskan apa yang dibutuhkan.
Selalu berikan insight yang actionable dan praktis.
Fokus pada membantu user membuat keputusan finansial yang lebih baik.

PENTING:
- Berikan jawaban yang spesifik dan berbasis data
- Hindari generalisasi tanpa data pendukung
- Selalu sertakan angka konkret jika tersedia
- Gunakan format rupiah (Rp) untuk nilai mata uang`

const AnalysisPromptTemplate = `User bertanya: "%s"

Data yang tersedia:
%s

Conversation history (5 pesan terakhir):
%s

Berikan analisis yang komprehensif mencakup:
1. Jawaban langsung untuk pertanyaan user (gunakan data yang spesifik)
2. Insights tambahan yang relevan dari data
3. Rekomendasi actionable (jika applicable)

Format response dalam JSON dengan struktur:
{
  "answer": "jawaban utama dengan data spesifik dan angka",
  "insights": ["insight 1 berbasis data", "insight 2 berbasis data"],
  "recommendations": ["rekomendasi 1 yang actionable", "rekomendasi 2 yang actionable"],
  "needs_visualization": true/false,
  "visualization_type": "bar/line/pie/table" (jika needs_visualization true),
  "visualization_hint": "penjelasan data apa yang perlu divisualisasikan"
}`

const RecommendationPromptTemplate = `Berdasarkan data finansial user berikut:

%s

Berikan rekomendasi untuk manajemen keuangan yang lebih baik.
Fokus pada area yang bisa dioptimasi dan berikan estimasi penghematan yang realistis.

Format response dalam JSON:
{
  "recommendations": [
    {
      "title": "judul singkat",
      "description": "penjelasan detail",
      "potential_saving": angka_dalam_rupiah,
      "difficulty": "easy/medium/hard",
      "action_items": ["langkah 1", "langkah 2"]
    }
  ],
  "total_potential_saving": total_estimasi_penghematan,
  "priority_order": ["recommendation yang paling penting dulu"]
}`

const VisualizationPromptTemplate = `User meminta visualisasi untuk query: "%s"

Data yang tersedia:
%s

Generate structure data untuk chart dengan format:
{
  "type": "bar/line/pie/table",
  "title": "judul chart",
  "data": {
    "labels": ["label1", "label2", ...],
    "values": [nilai1, nilai2, ...],
    "colors": ["#color1", "#color2", ...] (optional)
  },
  "config": {
    "x_label": "label sumbu x",
    "y_label": "label sumbu y",
    "format": "currency/percentage/number"
  }
}`

const TrendAnalysisPromptTemplate = `Analisis trend dari data berikut:

%s

Identifikasi:
1. Pola dan trend yang terlihat
2. Perubahan signifikan
3. Anomali atau outliers
4. Prediksi untuk periode berikutnya

Format response dalam JSON:
{
  "trend_direction": "increasing/decreasing/stable",
  "trend_percentage": persentase_perubahan,
  "key_insights": ["insight 1", "insight 2"],
  "predictions": {
    "next_period": nilai_prediksi,
    "confidence": "high/medium/low",
    "reasoning": "penjelasan prediksi"
  },
  "alerts": ["peringatan jika ada anomali"]
}`

const BudgetComparisonPromptTemplate = `Bandingkan budget vs actual spending:

Budget Data:
%s

Actual Spending:
%s

Analisis:
1. Kategori mana yang over/under budget
2. Tingkat adherence ke budget
3. Rekomendasi adjustment

Format response dalam JSON:
{
  "overall_status": "over_budget/under_budget/on_track",
  "variance_percentage": persentase_selisih,
  "category_breakdown": [
    {
      "category": "nama kategori",
      "budget": nilai_budget,
      "actual": nilai_actual,
      "variance": selisih,
      "variance_percentage": persentase,
      "status": "over/under/on_track"
    }
  ],
  "recommendations": ["rekomendasi 1", "rekomendasi 2"]
}`

const CategoryAnalysisPromptTemplate = `Analisis pengeluaran per kategori:

Data kategori:
%s

Periode: %s

Berikan breakdown dan insights:
1. Kategori dengan pengeluaran terbesar
2. Distribusi persentase
3. Perbandingan dengan periode sebelumnya (jika ada)
4. Rekomendasi optimasi

Format response dalam JSON:
{
  "top_categories": [
    {
      "category": "nama",
      "amount": nilai,
      "percentage": persentase_dari_total,
      "transaction_count": jumlah_transaksi
    }
  ],
  "insights": ["insight 1", "insight 2"],
  "optimization_opportunities": ["peluang 1", "peluang 2"]
}`

const SummaryPromptTemplate = `Generate ringkasan finansial untuk periode: %s

Data:
%s

Buat ringkasan komprehensif meliputi:
1. Total income dan expense
2. Net cashflow
3. Kategori terbesar
4. Perbandingan dengan periode sebelumnya
5. Highlight penting

Format dalam bahasa Indonesia yang mudah dipahami.`

// BuildAnalysisPrompt builds the prompt for general analysis
func BuildAnalysisPrompt(userQuery string, dbContext string, conversationHistory string) string {
	return fmt.Sprintf(AnalysisPromptTemplate, userQuery, dbContext, conversationHistory)
}

// BuildRecommendationPrompt builds the prompt for recommendations
func BuildRecommendationPrompt(userData string) string {
	return fmt.Sprintf(RecommendationPromptTemplate, userData)
}

// BuildVisualizationPrompt builds the prompt for visualization data
func BuildVisualizationPrompt(userQuery string, data string) string {
	return fmt.Sprintf(VisualizationPromptTemplate, userQuery, data)
}

// BuildTrendAnalysisPrompt builds the prompt for trend analysis
func BuildTrendAnalysisPrompt(trendData string) string {
	return fmt.Sprintf(TrendAnalysisPromptTemplate, trendData)
}

// BuildBudgetComparisonPrompt builds the prompt for budget comparison
func BuildBudgetComparisonPrompt(budgetData string, actualData string) string {
	return fmt.Sprintf(BudgetComparisonPromptTemplate, budgetData, actualData)
}

// BuildCategoryAnalysisPrompt builds the prompt for category analysis
func BuildCategoryAnalysisPrompt(categoryData string, period string) string {
	return fmt.Sprintf(CategoryAnalysisPromptTemplate, categoryData, period)
}

// BuildSummaryPrompt builds the prompt for summary generation
func BuildSummaryPrompt(period string, data string) string {
	return fmt.Sprintf(SummaryPromptTemplate, period, data)
}
