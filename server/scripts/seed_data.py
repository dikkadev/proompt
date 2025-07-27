#!/usr/bin/env -S uv run --script
# /// script
# dependencies = [
#   "faker",
# ]
# ///

"""
Seed the proompt database with sample data.
This script creates realistic test data for development and testing.

Edit the data generation functions below to customize what gets created:
- modify_sample_data(): Change the specific prompts, snippets, and notes
- SAMPLE_COUNTS: Adjust how many of each item to create
- generate_*(): Modify the data generation logic
"""

import sqlite3
import os
import sys
import json
import uuid
from datetime import datetime, timedelta
from faker import Faker

# Initialize Faker for generating realistic data
fake = Faker()

# Configuration - Edit these to change what gets generated
SAMPLE_COUNTS = {
    'prompts': 15,
    'snippets': 8,
    'notes_per_prompt': 2,  # Average notes per prompt
    'tags_per_item': 3,     # Average tags per prompt/snippet
    'prompt_links': 5,      # Number of prompt links to create
}

def get_db_path():
    """Get the database path from the config or use default."""
    return "./data/proompt.db"

def generate_prompt_data():
    """Generate sample prompt data. Edit this function to customize prompts."""
    prompts = []
    
    # Some predefined useful prompts with variable syntax and snippet references
    predefined_prompts = [
        {
            'title': 'Code Review Assistant',
            'content': 'You are an expert code reviewer for {{language:Python}}. Please review the following code and provide constructive feedback on:\n\n1. Code quality and best practices\n2. Potential bugs or issues\n3. Performance improvements\n4. Security considerations\n5. Readability and maintainability\n\n@error_handling_pattern\n\nCode to review:\n```{{language}}\n{{code}}\n```\n\nFocus areas: {{focus_areas:general review}}',
            'type': 'system',
            'use_case': 'code_review',
            'model_compatibility_tags': ['gpt-4', 'claude-3', 'gemini-pro'],
            'temperature_suggestion': 0.3,
        },
        {
            'title': 'Technical Documentation Writer',
            'content': 'You are a technical documentation specialist for {{project_name}}. Create clear, comprehensive documentation for:\n\n**Topic:** {{topic}}\n**Audience:** {{audience:developers}}\n**Format:** {{format:markdown}}\n\nInclude:\n- Overview and purpose\n- Prerequisites\n- Step-by-step instructions\n- Examples\n- Common troubleshooting\n- Best practices\n\n@api_response_format\n\nAdditional context: {{context}}',
            'type': 'system',
            'use_case': 'documentation',
            'model_compatibility_tags': ['gpt-4', 'claude-3'],
            'temperature_suggestion': 0.4,
        },
        {
            'title': 'API Design Consultant',
            'content': 'Design a RESTful API for {{domain}} with {{auth_type:JWT}} authentication.\n\nConsider:\n1. Resource modeling for {{resources}}\n2. HTTP methods and status codes\n3. Request/response formats\n4. Authentication and authorization\n5. Rate limiting ({{rate_limit:1000/hour}})\n6. Versioning strategy\n7. Error handling\n\n@input_validation\n@database_transaction\n\nProvide OpenAPI specification and implementation notes.\nTarget framework: {{framework:Express.js}}',
            'type': 'user',
            'use_case': 'api_design',
            'model_compatibility_tags': ['gpt-4', 'claude-3', 'gemini-pro'],
            'temperature_suggestion': 0.5,
        },
        {
            'title': 'Database Schema Designer',
            'content': 'Design a database schema for {{application_type}} using {{database:PostgreSQL}}.\n\nRequirements:\n- Expected users: {{user_count:10000}}\n- Data retention: {{retention:5 years}}\n- Performance target: {{performance:sub-100ms queries}}\n\nInclude:\n1. Entity relationship diagram\n2. Table definitions with constraints\n3. Indexes for performance\n4. Migration strategy\n5. Data integrity considerations\n\n@error_handling_pattern\n\nExplain your design decisions and trade-offs for {{specific_requirements}}.',
            'type': 'user',
            'use_case': 'database_design',
            'model_compatibility_tags': ['gpt-4', 'claude-3'],
            'temperature_suggestion': 0.4,
        },
        {
            'title': 'Security Audit Assistant',
            'content': 'Perform a security audit of {{system_name}} ({{system_type:web application}}).\n\nScope: {{audit_scope:full application}}\nCompliance: {{compliance_requirements:GDPR, SOC2}}\n\nFocus on:\n1. Authentication and authorization flaws\n2. Input validation issues\n3. SQL injection vulnerabilities\n4. XSS and CSRF risks\n5. Data exposure concerns\n6. Configuration security\n\n@input_validation\n@missing_snippet_reference\n\nSystem/Code:\n{{target}}\n\nPriority areas: {{priority_areas:authentication, data handling}}',
            'type': 'system',
            'use_case': 'security_audit',
            'model_compatibility_tags': ['gpt-4', 'claude-3'],
            'temperature_suggestion': 0.2,
        },
    ]
    
    # Add predefined prompts
    for prompt_data in predefined_prompts:
        prompt = {
            'id': str(uuid.uuid4()),
            'title': prompt_data['title'],
            'content': prompt_data['content'],
            'type': prompt_data['type'],
            'use_case': prompt_data['use_case'],
            'model_compatibility_tags': json.dumps(prompt_data['model_compatibility_tags']),
            'temperature_suggestion': prompt_data['temperature_suggestion'],
            'other_parameters': json.dumps({'max_tokens': 2000}),
            'created_at': fake.date_time_between(start_date='-30d', end_date='now').isoformat(),
            'updated_at': datetime.now().isoformat(),
            'git_ref': None,
        }
        prompts.append(prompt)
    
    # Generate additional random prompts with variable syntax and snippet references
    prompt_types = ['system', 'user', 'image', 'video']
    use_cases = ['development', 'testing', 'documentation', 'analysis', 'creative', 'research']
    
    # Templates with different variable syntax patterns
    content_templates = [
        "Analyze {{data_source}} and provide insights about {{topic:general trends}}. Focus on {{metrics}} and include @api_response_format in your analysis.",
        "Create a {{document_type:report}} for {{audience:stakeholders}} covering {{subject}}. Use {{format:markdown}} and reference @input_validation patterns.",
        "Generate {{content_type}} for {{platform:web}} targeting {{user_segment:general users}}. Include {{features}} and apply @error_handling_pattern.",
        "Review {{code_type:JavaScript}} code for {{project_name}} focusing on {{review_aspects:performance, security}}. Apply @database_transaction principles.",
        "Design {{system_component}} for {{application:web app}} with {{requirements}} and {{constraints:budget friendly}}. Reference @missing_snippet for advanced patterns.",
        "Implement {{feature_name}} using {{technology:React}} with {{styling:CSS modules}}. Consider {{accessibility_requirements:WCAG 2.1}} and use @api_response_format.",
        "Test {{functionality}} in {{environment:staging}} with {{test_data}} and {{expected_results:positive outcomes}}. Include @input_validation checks.",
        "Deploy {{service_name}} to {{platform:AWS}} with {{configuration}} and {{monitoring:CloudWatch}}. Follow @error_handling_pattern for failures.",
        "Optimize {{performance_target}} for {{system_part}} using {{optimization_method:caching}} and {{tools:profiler}}. Apply @database_transaction optimizations.",
        "Document {{api_endpoint}} with {{parameters}} returning {{response_format}} for {{use_case:user management}}. Use @api_response_format structure."
    ]
    
    for i in range(SAMPLE_COUNTS['prompts'] - len(predefined_prompts)):
        prompt_type = fake.random_element(prompt_types)
        use_case = fake.random_element(use_cases)
        content_template = fake.random_element(content_templates)
        
        prompt = {
            'id': str(uuid.uuid4()),
            'title': fake.catch_phrase(),
            'content': content_template,
            'type': prompt_type,
            'use_case': use_case,
            'model_compatibility_tags': json.dumps(fake.random_elements(['gpt-4', 'claude-3', 'gemini-pro', 'llama-2'], length=2)),
            'temperature_suggestion': round(fake.random.uniform(0.1, 1.0), 1),
            'other_parameters': json.dumps({'max_tokens': fake.random_int(500, 4000)}),
            'created_at': fake.date_time_between(start_date='-30d', end_date='now').isoformat(),
            'updated_at': datetime.now().isoformat(),
            'git_ref': None,
        }
        prompts.append(prompt)
    
    return prompts

def generate_snippet_data():
    """Generate sample snippet data. Edit this function to customize snippets."""
    snippets = []
    
    # Predefined useful snippets (these will be referenced by @ in prompts)
    predefined_snippets = [
        {
            'title': 'error_handling_pattern',
            'content': 'try {\n    // risky operation for {{operation_name:default operation}}\n    const result = await {{function_name}}({{params}});\n    return result;\n} catch (error) {\n    logger.error("{{error_message:Operation failed}}", error);\n    throw new CustomError("{{user_message:Operation failed}}", error);\n}',
            'description': 'Standard error handling pattern with logging and custom error wrapping',
        },
        {
            'title': 'database_transaction',
            'content': 'const transaction = await db.beginTransaction();\ntry {\n    await db.query("INSERT INTO {{table1:users}} ...", {{params1}});\n    await db.query("INSERT INTO {{table2:profiles}} ...", {{params2}});\n    await transaction.commit();\n    console.log("{{success_message:Transaction completed successfully}}");\n} catch (error) {\n    await transaction.rollback();\n    logger.error("Transaction failed for {{context}}", error);\n    throw error;\n}',
            'description': 'Database transaction pattern with rollback on error',
        },
        {
            'title': 'api_response_format',
            'content': '{\n  "success": {{success:true}},\n  "data": {{data}},\n  "message": "{{message:Operation completed successfully}}",\n  "timestamp": "{{timestamp}}",\n  "request_id": "{{request_id}}",\n  "metadata": {\n    "version": "{{api_version:v1}}",\n    "endpoint": "{{endpoint}}"\n  }\n}',
            'description': 'Standard API response format with metadata and variables',
        },
        {
            'title': 'input_validation',
            'content': 'const schema = {\n  email: { type: "string", format: "email", required: {{email_required:true}} },\n  age: { type: "number", minimum: {{min_age:0}}, maximum: {{max_age:150}} },\n  name: { type: "string", minLength: 1, maxLength: {{max_name_length:100}} },\n  {{custom_field}}: { type: "{{field_type:string}}", required: {{field_required:false}} }\n};\n\nconst isValid = validate(schema, {{input_data}});\nif (!isValid) {\n  throw new ValidationError("{{validation_message:Invalid input data}");\n}',
            'description': 'JSON schema validation pattern with configurable fields',
        },
    ]
    
    # Add predefined snippets
    for snippet_data in predefined_snippets:
        snippet = {
            'id': str(uuid.uuid4()),
            'title': snippet_data['title'],
            'content': snippet_data['content'],
            'description': snippet_data['description'],
            'created_at': fake.date_time_between(start_date='-30d', end_date='now').isoformat(),
            'updated_at': datetime.now().isoformat(),
            'git_ref': None,
        }
        snippets.append(snippet)
    
    # Generate additional random snippets with variable syntax
    snippet_templates = [
        {
            'title': 'auth_middleware',
            'content': 'function authenticate(req, res, next) {\n  const token = req.headers["{{auth_header:authorization}}"];\n  if (!token) {\n    return res.status(401).json({{error_response:{"error": "No token provided"}}});\n  }\n  \n  jwt.verify(token, {{secret_key}}, (err, decoded) => {\n    if (err) {\n      return res.status({{error_status:403}}).json({{invalid_token_response}});\n    }\n    req.user = decoded;\n    next();\n  });\n}',
            'description': 'JWT authentication middleware with configurable responses',
        },
        {
            'title': 'logging_config',
            'content': 'const logger = winston.createLogger({\n  level: "{{log_level:info}}",\n  format: winston.format.combine(\n    winston.format.timestamp(),\n    winston.format.errors({ stack: true }),\n    winston.format.json()\n  ),\n  defaultMeta: { service: "{{service_name}}" },\n  transports: [\n    new winston.transports.File({ filename: "{{error_log:error.log}}", level: "error" }),\n    new winston.transports.File({ filename: "{{combined_log:combined.log}}" })\n  ]\n});',
            'description': 'Winston logger configuration with customizable settings',
        },
        {
            'title': 'rate_limiter',
            'content': 'const rateLimit = require("express-rate-limit");\n\nconst limiter = rateLimit({\n  windowMs: {{window_minutes:15}} * 60 * 1000,\n  max: {{max_requests:100}},\n  message: "{{rate_limit_message:Too many requests from this IP}}",\n  standardHeaders: {{standard_headers:true}},\n  legacyHeaders: {{legacy_headers:false}},\n  handler: (req, res) => {\n    res.status({{rate_limit_status:429}}).json({\n      error: "{{custom_error_message:Rate limit exceeded}}",\n      retryAfter: {{retry_after:900}}\n    });\n  }\n});',
            'description': 'Express rate limiting middleware with configurable limits',
        },
        {
            'title': 'cache_helper',
            'content': 'class CacheHelper {\n  constructor() {\n    this.redis = new Redis({{redis_config}});\n    this.defaultTTL = {{default_ttl:3600}};\n  }\n\n  async get(key) {\n    const cached = await this.redis.get("{{key_prefix:app}}:" + key);\n    return cached ? JSON.parse(cached) : null;\n  }\n\n  async set(key, value, ttl = this.defaultTTL) {\n    await this.redis.setex(\n      "{{key_prefix:app}}:" + key,\n      ttl || {{fallback_ttl:1800}},\n      JSON.stringify(value)\n    );\n  }\n}',
            'description': 'Redis cache helper with configurable prefixes and TTL',
        },
    ]
    
    # Add the additional snippets
    for snippet_data in snippet_templates:
        snippet = {
            'id': str(uuid.uuid4()),
            'title': snippet_data['title'],
            'content': snippet_data['content'],
            'description': snippet_data['description'],
            'created_at': fake.date_time_between(start_date='-30d', end_date='now').isoformat(),
            'updated_at': datetime.now().isoformat(),
            'git_ref': None,
        }
        snippets.append(snippet)
    
    return snippets

def generate_notes_data(prompts):
    """Generate sample notes for prompts."""
    notes = []
    
    for prompt in prompts:
        # Generate 0-4 notes per prompt (average of 2)
        num_notes = fake.random_int(0, 4)
        
        for _ in range(num_notes):
            note = {
                'id': str(uuid.uuid4()),
                'prompt_id': prompt['id'],
                'title': fake.sentence(nb_words=4),
                'body': fake.text(max_nb_chars=200),
                'created_at': fake.date_time_between(start_date='-30d', end_date='now').isoformat(),
                'updated_at': datetime.now().isoformat(),
            }
            notes.append(note)
    
    return notes

def generate_tags_data(prompts, snippets):
    """Generate tags for prompts and snippets."""
    prompt_tags = []
    snippet_tags = []
    
    # Common tags
    common_tags = [
        'development', 'testing', 'documentation', 'api', 'database',
        'security', 'performance', 'frontend', 'backend', 'devops',
        'ai', 'ml', 'automation', 'review', 'analysis'
    ]
    
    # Generate prompt tags
    for prompt in prompts:
        num_tags = fake.random_int(1, 5)
        selected_tags = fake.random_elements(common_tags, length=num_tags, unique=True)
        
        for tag in selected_tags:
            prompt_tags.append({
                'prompt_id': prompt['id'],
                'tag_name': tag,
            })
    
    # Generate snippet tags
    for snippet in snippets:
        num_tags = fake.random_int(1, 4)
        selected_tags = fake.random_elements(common_tags, length=num_tags, unique=True)
        
        for tag in selected_tags:
            snippet_tags.append({
                'snippet_id': snippet['id'],
                'tag_name': tag,
            })
    
    return prompt_tags, snippet_tags

def generate_prompt_links(prompts):
    """Generate links between prompts."""
    links = []
    link_types = ['followup', 'related', 'prerequisite', 'alternative']
    
    for _ in range(min(SAMPLE_COUNTS['prompt_links'], len(prompts) - 1)):
        from_prompt = fake.random_element(prompts)
        to_prompt = fake.random_element([p for p in prompts if p['id'] != from_prompt['id']])
        
        link = {
            'from_prompt_id': from_prompt['id'],
            'to_prompt_id': to_prompt['id'],
            'link_type': fake.random_element(link_types),
            'created_at': fake.date_time_between(start_date='-30d', end_date='now').isoformat(),
        }
        links.append(link)
    
    return links

def seed_database():
    """Seed the database with sample data."""
    db_path = get_db_path()
    
    if not os.path.exists(db_path):
        print(f"Database file not found at: {db_path}")
        print("Make sure you're running this from the server directory and the database has been initialized.")
        sys.exit(1)
    
    print(f"Seeding database at: {db_path}")
    
    try:
        conn = sqlite3.connect(db_path)
        cursor = conn.cursor()
        
        # Generate data
        print("Generating sample data...")
        prompts = generate_prompt_data()
        snippets = generate_snippet_data()
        notes = generate_notes_data(prompts)
        prompt_tags, snippet_tags = generate_tags_data(prompts, snippets)
        prompt_links = generate_prompt_links(prompts)
        
        print(f"Generated:")
        print(f"  - {len(prompts)} prompts")
        print(f"  - {len(snippets)} snippets")
        print(f"  - {len(notes)} notes")
        print(f"  - {len(prompt_tags)} prompt tags")
        print(f"  - {len(snippet_tags)} snippet tags")
        print(f"  - {len(prompt_links)} prompt links")
        
        # Insert prompts
        print("Inserting prompts...")
        for prompt in prompts:
            cursor.execute("""
                INSERT INTO prompts (id, title, content, type, use_case, model_compatibility_tags, 
                                   temperature_suggestion, other_parameters, created_at, updated_at, git_ref)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """, (
                prompt['id'], prompt['title'], prompt['content'], prompt['type'],
                prompt['use_case'], prompt['model_compatibility_tags'],
                prompt['temperature_suggestion'], prompt['other_parameters'],
                prompt['created_at'], prompt['updated_at'], prompt['git_ref']
            ))
        
        # Insert snippets
        print("Inserting snippets...")
        for snippet in snippets:
            cursor.execute("""
                INSERT INTO snippets (id, title, content, description, created_at, updated_at, git_ref)
                VALUES (?, ?, ?, ?, ?, ?, ?)
            """, (
                snippet['id'], snippet['title'], snippet['content'], snippet['description'],
                snippet['created_at'], snippet['updated_at'], snippet['git_ref']
            ))
        
        # Insert notes
        print("Inserting notes...")
        for note in notes:
            cursor.execute("""
                INSERT INTO notes (id, prompt_id, title, body, created_at, updated_at)
                VALUES (?, ?, ?, ?, ?, ?)
            """, (
                note['id'], note['prompt_id'], note['title'], note['body'],
                note['created_at'], note['updated_at']
            ))
        
        # Insert prompt tags
        print("Inserting prompt tags...")
        for tag in prompt_tags:
            cursor.execute("""
                INSERT INTO prompt_tags (prompt_id, tag_name)
                VALUES (?, ?)
            """, (tag['prompt_id'], tag['tag_name']))
        
        # Insert snippet tags
        print("Inserting snippet tags...")
        for tag in snippet_tags:
            cursor.execute("""
                INSERT INTO snippet_tags (snippet_id, tag_name)
                VALUES (?, ?)
            """, (tag['snippet_id'], tag['tag_name']))
        
        # Insert prompt links
        print("Inserting prompt links...")
        for link in prompt_links:
            cursor.execute("""
                INSERT INTO prompt_links (from_prompt_id, to_prompt_id, link_type, created_at)
                VALUES (?, ?, ?, ?)
            """, (link['from_prompt_id'], link['to_prompt_id'], link['link_type'], link['created_at']))
        
        # Commit changes
        conn.commit()
        print("\n✅ Database seeded successfully!")
        
    except sqlite3.Error as e:
        print(f"Database error: {e}")
        sys.exit(1)
    finally:
        if conn:
            conn.close()

def modify_sample_data():
    """
    Modify this function to customize the sample data generation.
    You can:
    - Change SAMPLE_COUNTS to generate more/fewer items
    - Edit the predefined prompts and snippets
    - Modify the data generation logic
    - Add new data types or relationships
    """
    print("💡 To customize the data generation:")
    print("   - Edit SAMPLE_COUNTS at the top of this file")
    print("   - Modify generate_prompt_data() and generate_snippet_data()")
    print("   - Customize the predefined prompts and snippets")
    print("   - Adjust the tags and relationships")

if __name__ == "__main__":
    print("🌱 Proompt Database Seeder")
    print("=" * 30)
    
    modify_sample_data()
    print()
    
    # Confirm action
    response = input("Proceed with seeding the database? (yes/no): ")
    if response.lower() not in ['yes', 'y']:
        print("Operation cancelled.")
        sys.exit(0)
    
    seed_database()