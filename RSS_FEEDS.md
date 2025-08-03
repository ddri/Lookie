# Working RSS Feeds for Quantum Computing

## ✅ Verified Working Feeds

### Government & Research Labs
- **Fermilab Quantum**: `https://news.fnal.gov/tag/quantum-computing/feed/`
  - 10 articles about quantum computing
  - Mentions companies like Diraq, HRL Labs
  - High-quality content about industry developments

### General Tech News (Quantum Coverage)
- **TechCrunch**: `https://techcrunch.com/feed/`
  - Covers quantum startup funding and news
  - 20+ articles available
  - Good source for business developments

### Quantum Computing Publications  
- **Quantum Computing Report**: `https://quantumcomputingreport.com/feed/`
  - Dedicated quantum computing news
  - Industry analysis and company updates
  - Regular coverage of major players

## ❌ Broken Company Feeds

### Quantum Companies (Need Alternative Sources)
- **IonQ**: `https://ionq.com/blog/feed` ❌ (Returns HTML, not RSS)
- **Q-CTRL**: `https://q-ctrl.com/blog/feed` ❌ (404 Not Found)
- **PsiQuantum**: No RSS feed found
- **Diraq**: No RSS feed found

## 🔧 Recommended Configuration

For immediate testing, update company RSS URLs to:

1. **IonQ** → `https://news.fnal.gov/tag/quantum-computing/feed/`
   - Will capture articles mentioning IonQ
   - Government research perspective

2. **Q-CTRL** → `https://quantumcomputingreport.com/feed/`
   - Industry news covering all quantum companies
   - Business and technical updates

3. **PsiQuantum** → `https://techcrunch.com/feed/`
   - Startup and funding news
   - Covers quantum industry developments

4. **Diraq** → `https://news.fnal.gov/tag/quantum-computing/feed/`
   - Already mentions Diraq in recent articles
   - Research collaborations and partnerships

## 🎯 How to Use Admin Panel

1. **Start Admin Panel**: `go run cmd/admin/main.go`
2. **Open Dashboard**: http://localhost:3000/admin.html
3. **Update RSS URLs**: Use the interface to set working feeds
4. **Test Feeds**: Click "Test" button for each company
5. **Trigger Scrape**: Use "Trigger Scrape" to collect articles
6. **Monitor Results**: Check dashboard stats

## 📋 Next Steps

1. Use admin panel to configure working RSS feeds
2. Test the complete pipeline: RSS → Firestore → Gemini AI
3. Build user dashboard for consuming intelligence
4. Add Cloud Scheduler for automated scraping

## 🔍 Finding More RSS Feeds

### Search Strategies:
- Check `/feed`, `/rss`, `/atom` on company domains
- Look for RSS icons on news pages
- Use feed discovery tools
- Monitor quantum industry publications

### Alternative Approaches (Phase 2):
- Direct website scraping of news pages
- Social media monitoring (Twitter, LinkedIn)
- Academic paper tracking (arXiv feeds)
- Conference and event monitoring