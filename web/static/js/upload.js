const fileInput = document.getElementById('fileInput');
const fileList = document.getElementById('fileList');
const folderNameInput = document.getElementById('folderName');
const uploadBtn = document.getElementById('uploadBtn');
const messageDiv = document.getElementById('message');

let selectedFiles = [];

fileInput.addEventListener('change', (e) => {
    selectedFiles = Array.from(e.target.files);
    displayFileList();
});

function displayFileList() {
    fileList.innerHTML = '';
    if (selectedFiles.length === 0) {
        return;
    }

    selectedFiles.forEach((file, index) => {
        const fileItem = document.createElement('div');
        fileItem.className = 'file-item';
        fileItem.innerHTML = `
            <div>
                <span>${file.name}</span>
                <small>(${formatFileSize(file.size)})</small>
            </div>
        `;
        fileList.appendChild(fileItem);
    });
}

function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

uploadBtn.addEventListener('click', async () => {
    const folderName = folderNameInput.value.trim();

    if (!folderName) {
        showMessage('Please enter a folder name', 'error');
        return;
    }

    if (selectedFiles.length === 0) {
        showMessage('Please select at least one file', 'error');
        return;
    }

    // Disable button during upload
    uploadBtn.disabled = true;
    uploadBtn.textContent = 'Uploading...';

    try {
        // Create folder first
        const folderResponse = await fetch('/api/folders', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name: folderName })
        });

        if (!folderResponse.ok) {
            throw new Error('Failed to create folder');
        }

        const folder = await folderResponse.json();

        // Upload files
        for (const file of selectedFiles) {
            const formData = new FormData();
            formData.append('file', file);
            formData.append('folder_id', folder.id);
            formData.append('folder_name', folder.name);

            const response = await fetch('/api/upload', {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                throw new Error(`Failed to upload ${file.name}`);
            }
        }

        showMessage(`Successfully uploaded ${selectedFiles.length} file(s) to folder "${folderName}"!`, 'success');
        
        // Reset form
        folderNameInput.value = '';
        fileInput.value = '';
        selectedFiles = [];
        fileList.innerHTML = '';

    } catch (error) {
        showMessage(`Error: ${error.message}`, 'error');
    } finally {
        uploadBtn.disabled = false;
        uploadBtn.textContent = 'Upload Files';
    }
});

function showMessage(text, type) {
    messageDiv.textContent = text;
    messageDiv.className = `message ${type}`;
    messageDiv.style.display = 'block';

    setTimeout(() => {
        messageDiv.style.display = 'none';
    }, 5000);
}
