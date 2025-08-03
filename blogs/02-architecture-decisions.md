# Building a Real-time Intelligence Service: Architecture Decisions for Lookie

Designing a production intelligence service requires careful consideration of scalability, reliability, and cost efficiency. Lookie's architecture reflects decisions made to optimize for the specific characteristics of news monitoring workloads while maintaining operational simplicity and development velocity.

## Document Store vs Relational Database Trade-offs

The choice between document-oriented and relational storage significantly impacts system design for content aggregation services. Lookie initially implemented a SQLite-based approach with normalized tables for companies, articles, and classifications. This design provided strong consistency and familiar query patterns but created performance bottlenecks as content volume increased.

The migration to Google Firestore addressed several fundamental limitations of the relational approach. News articles exhibit variable schema characteristics that map poorly to fixed table structures. Article metadata varies significantly between sources, with some providing detailed author information, publication dates, and categories while others offer minimal structured data.

Firestore's document model accommodates this variability naturally. Articles are stored as documents with embedded metadata, allowing each document to contain different field sets without schema migration overhead. The system can capture rich metadata when available while gracefully handling minimal data from limited sources.

Query performance considerations also favored the document approach. Common access patterns in intelligence services involve retrieving articles with their associated company context and classification results. The relational model required join operations across multiple tables, creating query complexity and performance overhead. Firestore's document embedding capabilities enable complete context retrieval through single document reads.

The document model also simplified real-time subscription implementation. Firestore provides native support for document change listeners, enabling immediate notification when new articles are classified or company profiles are updated. Implementing equivalent functionality with relational databases requires additional message queuing infrastructure and complex change detection logic.

However, the document approach introduces trade-offs in data consistency and analytical querying. Firestore's eventual consistency model requires careful consideration of read-after-write scenarios. The system implements client-side retry logic for cases where immediately-written data may not be visible to subsequent reads.

Complex analytical queries that span multiple document types require careful index design and sometimes necessitate data denormalization. Lookie addresses this through strategic embedding of frequently-accessed relationship data within documents, trading storage efficiency for query performance.

## Event-driven Processing Patterns

Lookie implements an event-driven architecture that decouples content acquisition from processing while maintaining system responsiveness. This pattern addresses the variable processing time characteristics of AI classification workloads and enables independent scaling of system components.

Content acquisition operates on scheduled intervals using Google Cloud Scheduler. This approach provides predictable resource utilization and respects source system rate limits. When new articles are discovered, they are immediately persisted to Firestore with a processing status indicating pending classification.

Document creation events trigger background processing through Firestore triggers connected to Cloud Functions. This pattern ensures that classification processing begins immediately upon content acquisition without requiring polling mechanisms or batch processing delays.

The event-driven approach enables graceful handling of processing failures. If AI classification fails due to service outages or rate limiting, the article remains in pending status until processing can be retried. The system implements exponential backoff retry logic to handle transient failures without overwhelming downstream services.

Processing results are written back to the original article document, triggering additional events for notification and indexing workflows. This pattern maintains loose coupling between processing stages while ensuring that state changes propagate through the system reliably.

The event architecture also supports horizontal scaling of processing components. Multiple Cloud Function instances can process classification requests concurrently, with Firestore providing natural work distribution through document-level locking. Processing capacity scales automatically based on workload without requiring explicit resource management.

## Scaling Considerations for Enterprise Deployment

Enterprise deployment scenarios require careful attention to scaling characteristics across multiple dimensions. Lookie's architecture addresses scaling challenges through component isolation and resource optimization strategies.

Content acquisition scaling differs fundamentally from processing scaling. RSS feed monitoring involves periodic bulk operations with predictable resource requirements. Processing scaling responds to variable AI service latency and content classification complexity. The system isolates these concerns through separate deployment units that scale independently.

Storage scaling leverages Firestore's automatic sharding and replication capabilities. Document design considers key distribution to avoid hotspots as content volume increases. Article documents use compound keys that include company identifiers and timestamp components, ensuring even distribution across storage nodes.

API scaling utilizes Google Cloud Run's automatic scaling capabilities with carefully tuned concurrency limits. The service implements connection pooling for Firestore clients and request batching for bulk operations. Response caching reduces database load for frequently-accessed content while maintaining data freshness through cache invalidation on content updates.

Network scaling addresses the challenge of monitoring increasing numbers of content sources without overwhelming external systems. The service implements distributed rate limiting that respects individual source limitations while maintaining overall system throughput. Circuit breaker patterns prevent cascade failures when individual sources become unavailable.

Cost scaling optimization balances performance requirements with operational expenses. The serverless architecture ensures that resource costs align with actual utilization rather than peak capacity provisioning. Reserved capacity options provide cost optimization for predictable baseline workloads while maintaining burst capacity for processing spikes.

## State Management and Consistency Patterns

Managing state consistency in distributed intelligence services requires careful consideration of data flow patterns and failure scenarios. Lookie implements several consistency patterns that balance reliability with performance requirements.

Content acquisition implements idempotency patterns that handle duplicate detection across multiple scraping cycles. The system uses content hashing to identify identical articles from different sources or time periods. Hash-based deduplication operates at the document level, preventing storage of duplicate content while maintaining source attribution metadata.

Processing state management tracks article lifecycle through explicit status fields rather than implicit state inference. Articles progress through defined states: discovered, processing, classified, and error conditions. State transitions are atomic operations that maintain consistency even during partial failure scenarios.

Classification result consistency addresses the challenge of updating document embeddings when AI processing completes. The system implements optimistic concurrency control to prevent conflicting updates to article documents. Retry logic handles conflicts by reloading current document state before applying classification results.

Cross-document consistency for aggregated statistics uses eventually consistent patterns appropriate for analytical workloads. Company-level statistics are updated through background processes that can tolerate temporary inconsistency. The system provides both real-time individual document access and eventually consistent aggregate views.

Transaction boundaries are carefully designed to minimize distributed transaction requirements while maintaining data integrity. Related updates are batched within single document transactions when possible. Cross-document updates that require strong consistency use Firestore transaction groups with appropriate conflict detection and retry logic.

The consistency model acknowledges that perfect consistency is unnecessary for many intelligence service use cases. Users can tolerate slight delays in aggregate statistics or classification updates as long as individual document integrity is maintained. This recognition enables performance optimizations that would be impossible under strict consistency requirements.

## Integration Architecture and API Design

External system integration represents a critical architecture consideration for intelligence services deployed in enterprise environments. Lookie's integration architecture provides multiple interaction patterns while maintaining security and operational control.

The REST API implements resource-oriented design patterns that align with standard HTTP semantics. Article resources support standard CRUD operations with additional endpoints for classification management and search functionality. The API design emphasizes predictable behavior and comprehensive error handling to support reliable integration development.

Authentication and authorization use Google Cloud Identity and Access Management integration for enterprise deployment scenarios. The system supports both service account authentication for automated integrations and user-based authentication for interactive applications. API rate limiting prevents abuse while accommodating legitimate high-volume usage patterns.

Webhook integration enables real-time notification delivery to external systems without requiring polling. The system implements reliable delivery patterns with retry logic and dead letter queues for handling temporary subscriber outages. Webhook payloads include sufficient context to enable processing without additional API calls when possible.

Batch export capabilities support integration with data warehouse and business intelligence platforms. The system provides both streaming export for real-time integration and bulk export for historical analysis. Export formats include both JSON for flexible processing and structured formats optimized for analytical workloads.

The integration architecture maintains backwards compatibility through API versioning strategies that enable incremental feature adoption without breaking existing integrations. Version negotiation allows clients to specify required functionality while enabling server-side optimization for different client capabilities.

## Operational Excellence and Monitoring

Production intelligence services require comprehensive observability and operational control capabilities. Lookie implements monitoring and management patterns that provide visibility into system behavior while enabling proactive issue resolution.

Structured logging provides detailed insight into system operation without overwhelming log volume. Log entries include correlation identifiers that enable request tracing across distributed components. Log sampling reduces storage costs while maintaining sufficient detail for troubleshooting and performance analysis.

Health check implementation goes beyond simple service availability to include dependency health and business logic validation. Health endpoints validate Firestore connectivity, AI service availability, and data freshness indicators. This comprehensive health model enables intelligent load balancing and automatic failure recovery.

Performance monitoring tracks both technical metrics and business-relevant indicators. Technical metrics include response times, error rates, and resource utilization. Business metrics track content acquisition rates, classification accuracy, and processing latency. Combined monitoring provides complete visibility into system effectiveness.

Alerting strategies balance notification completeness with operational noise reduction. Critical alerts trigger immediate notification for service outages or data corruption scenarios. Warning alerts aggregate over time periods to identify trends without generating excessive notifications. The alerting system includes context-rich information that enables rapid issue diagnosis and resolution.

Deployment automation ensures consistent and reliable release processes. The system implements blue-green deployment patterns that enable zero-downtime updates with automatic rollback capabilities. Deployment includes comprehensive validation steps that verify system functionality before completing traffic migration.

This architecture demonstrates how thoughtful design decisions can create production-ready intelligence services that scale effectively while maintaining operational simplicity. The patterns and practices implemented in Lookie provide a foundation for similar systems facing comparable challenges in content processing and real-time intelligence delivery.