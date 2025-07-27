# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a 3D printing farm system (云3D打印农场系统) project in early planning/design phase. The repository currently contains only design documentation describing a cloud-based 3D printing service platform.

## Architecture Design

Based on README.MD, the planned system architecture includes:

### Frontend
- Vue.js framework
- Three.js for 3D model visualization
- Responsive design for desktop and mobile

### Backend  
- Go language framework
- REST API architecture
- File storage system for 3D models

### Database
- PostgreSQL for user data and orders
- File system storage for 3D model files

## Core Features (Planned)

1. **User Management**: Registration/login with email and password
2. **3D Model Processing**: STL/OBJ file upload (<50MB) with web-based 3D preview
3. **Printing Services**: Material selection (PLA/ABS/PETG), quality presets, automatic pricing
4. **Order Management**: Shopping cart interface, payment integration, order tracking

## Development Status

**Note**: This repository appears to be in the initial planning phase with only documentation present. No source code has been implemented yet.

## File Structure

```
/
├── README.MD           # Project design specification in Chinese
└── .claude/           # Claude Code configuration
    └── settings.local.json
```

When implementing this project, consider:
- Following the planned tech stack (Vue.js + Go + PostgreSQL)
- Implementing the 4-tier order status system described in the design
- Supporting the specified 3D model formats and file size limits
- Building the material selection and pricing system as outlined