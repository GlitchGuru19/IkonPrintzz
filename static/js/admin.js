/**
 * Admin Dashboard - Real-time file management with WebSocket updates
 * Handles file viewing, printing, downloading, and deletion
 */
class AdminDashboard {
    constructor() {
        this.files = [];
        this.ws = null;
        this.isConnected = false;
        
        this.initElements();
        this.initEventListeners();
        this.connectWebSocket();
        this.loadFiles();
    }

    /**
     * Initialize DOM elements
     */
    initElements() {
        this.filesTable = document.getElementById('filesTable');
        this.filesTableBody = document.getElementById('filesTableBody');
        this.loading = document.getElementById('loading');
        this.emptyState = document.getElementById('emptyState');
        this.fileCount = document.getElementById('fileCount');
        this.refreshBtn = document.getElementById('refreshBtn');
        this.cleanPrintedBtn = document.getElementById('cleanPrintedBtn');
    }

    /**
     * Initialize event listeners
     */
    initEventListeners() {
        this.refreshBtn.addEventListener('click', () => this.loadFiles());
        this.cleanPrintedBtn.addEventListener('click', () => this.cleanPrintedFiles());
    }

    /**
     * Connect to WebSocket for real-time updates
     */
    connectWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws/admin`;
        
        this.ws = new WebSocket(wsUrl);
        
        this.ws.onopen = () => {
            console.log('WebSocket connected - receiving live updates');
            this.isConnected = true;
            this.updateConnectionStatus();
        };
        
        this.ws.onclose = () => {
            console.log('WebSocket disconnected');
            this.isConnected = false;
            this.updateConnectionStatus();
            
            // Attempt to reconnect after 5 seconds
            setTimeout(() => this.connectWebSocket(), 5000);
        };
        
        this.ws.onmessage = (event) => {
            const files = JSON.parse(event.data);
            this.updateFiles(files);
        };
        
        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    /**
     * Update connection status indicator
     */
    updateConnectionStatus() {
        // Remove existing status indicator
        const existingStatus = document.querySelector('.connection-status');
        if (existingStatus) {
            existingStatus.remove();
        }

        const status = document.createElement('div');
        status.className = `connection-status ${this.isConnected ? 'connected' : 'disconnected'}`;
        status.innerHTML = this.isConnected ? 
            '<i class="fas fa-wifi me-1"></i> Live Updates' :
            '<i class="fas fa-wifi-slash me-1"></i> Reconnecting...';
        
        document.body.appendChild(status);
    }

    /**
     * Load files via HTTP API (fallback if WebSocket fails)
     */
    async loadFiles() {
        try {
            const response = await fetch('/files');
            const files = await response.json();
            this.updateFiles(files);
        } catch (error) {
            console.error('Error loading files:', error);
        }
    }

    /**
     * Update files list and re-render
     */
    updateFiles(files) {
        this.files = files;
        this.renderFiles();
        this.updateUI();
    }

    /**
     * Render files table
     */
    renderFiles() {
        if (this.files.length === 0) {
            this.filesTable.classList.add('d-none');
            this.emptyState.classList.remove('d-none');
            this.loading.classList.add('d-none');
            return;
        }

        this.emptyState.classList.add('d-none');
        this.loading.classList.add('d-none');
        this.filesTable.classList.remove('d-none');

        this.filesTableBody.innerHTML = this.files.map(file => `
            <tr class="file-row ${file.is_processed ? 'printed' : ''}">
                <td>
                    <code class="fs-6">${file.upload_code}</code>
                </td>
                <td>
                    <strong>${file.original_name}</strong>
                </td>
                <td>${this.formatFileSize(file.file_size)}</td>
                <td>${this.formatDate(file.upload_date)}</td>
                <td>
                    <span class="badge bg-secondary">${this.getFileType(file.content_type)}</span>
                </td>
                <td>
                    ${file.is_processed ? 
                        '<span class="badge bg-success"><i class="fas fa-check me-1"></i>Printed</span>' : 
                        '<span class="badge bg-warning"><i class="fas fa-clock me-1"></i>Pending</span>'}
                </td>
                <td>
                    <div class="btn-group btn-group-sm">
                        <button class="btn btn-outline-primary" onclick="admin.viewFile('${file.id}')" 
                                title="View File">
                            <i class="fas fa-eye"></i>
                        </button>
                        <button class="btn btn-outline-success" onclick="admin.printFile('${file.id}')" 
                                ${file.is_processed ? 'disabled' : ''}
                                title="Print File">
                            <i class="fas fa-print"></i>
                        </button>
                        <button class="btn btn-outline-info" onclick="admin.downloadFile('${file.id}')" 
                                title="Download">
                            <i class="fas fa-download"></i>
                        </button>
                        <button class="btn btn-outline-danger" onclick="admin.deleteFile('${file.id}')" 
                                title="Delete">
                            <i class="fas fa-trash"></i>
                        </button>
                    </div>
                </td>
            </tr>
        `).join('');
    }

    /**
     * Update UI elements based on current state
     */
    updateUI() {
        this.fileCount.textContent = this.files.length;
        
        const hasPrintedFiles = this.files.some(file => file.is_processed);
        this.cleanPrintedBtn.disabled = !hasPrintedFiles;
    }

    /**
     * Print a file (opens in new tab for printing)
     */
    async printFile(fileId) {
        const file = this.files.find(f => f.id === fileId);
        if (!file) return;

        try {
            // Open file in new tab for printing
            const fileUrl = `/uploads/${file.file_name}`;
            const printWindow = window.open(fileUrl, '_blank');
            
            if (printWindow) {
                printWindow.onload = () => {
                    printWindow.print();
                };
            }

            // Mark as printed on server
            await fetch(`/print/${fileId}`, { method: 'POST' });
            
        } catch (error) {
            console.error('Error printing file:', error);
            alert('Error printing file: ' + error.message);
        }
    }

    /**
     * View file in new tab
     */
    viewFile(fileId) {
        const file = this.files.find(f => f.id === fileId);
        if (file) {
            const fileUrl = `/uploads/${file.file_name}`;
            window.open(fileUrl, '_blank');
        }
    }

    /**
     * Download file to local device
     */
    downloadFile(fileId) {
        const file = this.files.find(f => f.id === fileId);
        if (file) {
            const fileUrl = `/uploads/${file.file_name}`;
            const a = document.createElement('a');
            a.href = fileUrl;
            a.download = file.original_name;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
        }
    }

    /**
     * Delete a file from server
     */
    async deleteFile(fileId) {
        if (!confirm('Are you sure you want to delete this file?')) {
            return;
        }

        try {
            const response = await fetch(`/delete/${fileId}`, { method: 'DELETE' });
            const result = await response.json();
            
            if (result.success) {
                this.loadFiles(); // Reload to get updated list
            } else {
                alert('Error deleting file: ' + result.message);
            }
        } catch (error) {
            console.error('Error deleting file:', error);
            alert('Error deleting file: ' + error.message);
        }
    }

    /**
     * Clean up all printed files
     */
    async cleanPrintedFiles() {
        const printedFiles = this.files.filter(file => file.is_processed);
        if (printedFiles.length === 0) return;

        if (!confirm(`Are you sure you want to delete all ${printedFiles.length} printed files?`)) {
            return;
        }

        try {
            for (const file of printedFiles) {
                await fetch(`/delete/${file.id}`, { method: 'DELETE' });
            }
            this.loadFiles(); // Reload to get updated list
        } catch (error) {
            console.error('Error cleaning printed files:', error);
            alert('Error cleaning printed files: ' + error.message);
        }
    }

    /**
     * Get file type from content type
     */
    getFileType(contentType) {
        if (contentType.startsWith('image/')) return 'Image';
        if (contentType === 'application/pdf') return 'PDF';
        if (contentType.includes('word')) return 'Word';
        if (contentType.includes('excel') || contentType.includes('sheet')) return 'Excel';
        if (contentType.includes('powerpoint') || contentType.includes('presentation')) return 'PowerPoint';
        if (contentType.includes('text')) return 'Text';
        return 'Document';
    }

    /**
     * Format file size for display
     */
    formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    /**
     * Format date for display
     */
    formatDate(dateString) {
        const date = new Date(dateString);
        return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { 
            hour: '2-digit', 
            minute: '2-digit' 
        });
    }
}

// Initialize admin dashboard
let admin;
document.addEventListener('DOMContentLoaded', () => {
    admin = new AdminDashboard();
});