// 全局变量
let currentStats = {};

// 分类配置
const CATEGORY_CONFIG = {
    '合同': {
        icon: 'fas fa-file-contract',
        color: 'work',
        description: '包含合同协议等相关文件'
    },
    '简历': {
        icon: 'fas fa-user-tie',
        color: 'personal',
        description: '包含个人简历和求职相关文件'
    },
    '发票': {
        icon: 'fas fa-receipt',
        color: 'images',
        description: '包含发票收据等财务文件'
    },
    '论文': {
        icon: 'fas fa-graduation-cap',
        color: 'study',
        description: '包含学术论文和研究文件'
    },
    '未分类': {
        icon: 'fas fa-question-circle',
        color: 'unclassified',
        description: '暂未识别类型的文件'
    },
    '新增分类': {
        icon: 'fas fa-plus-circle',
        color: 'new',
        description: '自定义分类，可自动分析'
    }
};

// DOM 元素
const uploadArea = document.getElementById('uploadArea');
const fileInput = document.getElementById('fileInput');
const uploadProgress = document.getElementById('uploadProgress');
const progressFill = document.getElementById('progressFill');
const progressText = document.getElementById('progressText');
const categoriesGrid = document.getElementById('categoriesGrid');
const fileListModal = document.getElementById('fileListModal');
const modalOverlay = document.getElementById('modalOverlay');

// 初始化
document.addEventListener('DOMContentLoaded', function() {
    initializeUpload();
    autoScanUploads();
});

// 自动扫描uploads文件夹
async function autoScanUploads() {
    try {
        console.log('开始自动扫描uploads文件夹...');
        
        const response = await fetch('/api/scan-uploads', {
            method: 'POST'
        });
        
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}`);
        }
        
        const result = await response.json();
        
        if (result.success) {
            console.log('uploads文件夹扫描完成:', result);
            if (result.results && result.results.total > 0) {
                updateSubtitle(result.results.total);
            }
        } else {
            console.error('扫描失败:', result.error);
        }
        
        // 加载统计数据
        loadStats();
        
    } catch (error) {
        console.error('自动扫描失败:', error);
        // 即使扫描失败，也要加载统计数据
        loadStats();
    }
}

// 初始化上传功能
function initializeUpload() {
    // 点击上传区域
    uploadArea.addEventListener('click', () => {
        fileInput.click();
    });

    // 文件选择
    fileInput.addEventListener('change', handleFileSelect);

    // 拖拽功能
    uploadArea.addEventListener('dragover', (e) => {
        e.preventDefault();
        uploadArea.classList.add('dragover');
    });

    uploadArea.addEventListener('dragleave', (e) => {
        e.preventDefault();
        uploadArea.classList.remove('dragover');
    });

    uploadArea.addEventListener('drop', (e) => {
        e.preventDefault();
        uploadArea.classList.remove('dragover');
        
        const files = Array.from(e.dataTransfer.files);
        if (files.length > 0) {
            uploadFiles(files);
        }
    });
}

// 处理文件选择
function handleFileSelect(event) {
    const files = Array.from(event.target.files);
    if (files.length > 0) {
        uploadFiles(files);
    }
}

// 上传文件
async function uploadFiles(files) {
    if (files.length === 0) {
        alert('请选择文件');
        return;
    }

    if (files.length > 200) {
        alert('最多支持上传200个文件');
        return;
    }

    // 显示进度条
    uploadArea.style.display = 'none';
    uploadProgress.style.display = 'block';
    
    const formData = new FormData();
    
    files.forEach(file => {
        formData.append('files', file);
        
        // 尝试获取文件路径信息
        let originalPath = file.name; // 默认使用文件名
        
        // 如果是文件夹上传，使用webkitRelativePath
        if (file.webkitRelativePath) {
            originalPath = file.webkitRelativePath;
            console.log('文件夹上传 - 相对路径:', file.webkitRelativePath);
        }
        // 如果是单个文件，尝试其他方法
        else {
            console.log('单个文件上传 - 文件名:', file.name);
            
            // 尝试获取更多文件信息
            const fileInfo = {
                name: file.name,
                size: file.size,
                type: file.type,
                lastModified: file.lastModified,
                webkitRelativePath: file.webkitRelativePath,
                path: file.path, // 在某些浏览器中可能可用
                mozFullPath: file.mozFullPath // Firefox中可能可用
            };
            
            console.log('文件详细信息:', fileInfo);
            
            // 如果有path属性，使用它
            if (file.path) {
                originalPath = file.path;
                console.log('使用path属性:', file.path);
            }
            // 如果有mozFullPath，使用它
            else if (file.mozFullPath) {
                originalPath = file.mozFullPath;
                console.log('使用mozFullPath:', file.mozFullPath);
            }
            // 否则使用文件名
            else {
                originalPath = file.name;
                console.log('使用文件名:', file.name);
            }
        }
        
        formData.append('originalPaths', originalPath);
    });

    try {
        // 模拟上传进度
        let progress = 0;
        const progressInterval = setInterval(() => {
            progress += Math.random() * 15;
            if (progress > 90) progress = 90;
            updateProgress(progress, `正在处理文件... ${Math.floor(progress)}%`);
        }, 300);

        const response = await fetch('/upload', {
            method: 'POST',
            body: formData
        });

        clearInterval(progressInterval);

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}`);
        }

        const result = await response.json();
        
        // 完成进度
        updateProgress(100, '文件分类完成！');
        
        setTimeout(() => {
            uploadProgress.style.display = 'none';
            uploadArea.style.display = 'block';
            loadStats();
            updateSubtitle(result.results.total);
        }, 1000);

        console.log('上传成功:', result);
        
    } catch (error) {
        console.error('上传失败:', error);
        alert('上传失败: ' + error.message);
        
        uploadProgress.style.display = 'none';
        uploadArea.style.display = 'block';
    }
}

// 更新进度条
function updateProgress(percent, text) {
    progressFill.style.width = `${percent}%`;
    progressText.textContent = text;
}

// 更新副标题
function updateSubtitle(count) {
    const subtitle = document.querySelector('.subtitle');
    subtitle.textContent = `已自动分类文件${count}个`;
}

// 加载统计数据
async function loadStats() {
    try {
        const response = await fetch('/api/stats');
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}`);
        }
        
        currentStats = await response.json();
        renderCategories();
        
    } catch (error) {
        console.error('加载统计数据失败:', error);
        // 显示默认的空状态
        renderCategories();
    }
}

// 渲染分类卡片
function renderCategories() {
    categoriesGrid.innerHTML = '';
    
    Object.entries(CATEGORY_CONFIG).forEach(([categoryName, config]) => {
        const stats = currentStats[categoryName] || { count: 0, files: [] };
        
        const card = document.createElement('div');
        card.className = `category-card ${config.color}`;
        card.innerHTML = `
            <i class="${config.icon} category-icon"></i>
            <h3 class="category-title">${categoryName}</h3>
            <div class="category-count">${stats.count}</div>
            <p class="category-description">${config.description}</p>
        `;
        
        // 添加点击事件
        card.addEventListener('click', () => {
            showFileList(categoryName, stats);
        });
        
        categoriesGrid.appendChild(card);
    });
}

// 显示文件列表
async function showFileList(categoryName, stats) {
    if (stats.count === 0) {
        alert(`${categoryName} 分类中暂无文件`);
        return;
    }

    try {
        // 获取最新的文件列表
        const response = await fetch(`/api/files/${encodeURIComponent(categoryName)}`);
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}`);
        }
        
        const data = await response.json();
        
        document.getElementById('modalTitle').textContent = `${categoryName} (${data.count}个文件)`;
        
        const fileList = document.getElementById('fileList');
        fileList.innerHTML = '';
        
        if (data.files && data.files.length > 0) {
            data.files.forEach(file => {
                const fileItem = document.createElement('div');
                fileItem.className = 'file-item';
                fileItem.style.cursor = 'pointer';
                
                const fileIcon = getFileIcon(file.name);
                const fileSize = formatFileSize(file.size);
                const badge = file.type === 'ai' ? '<span class="file-badge ai">AI分析</span>' : 
                             file.type === 'filename' ? '<span class="file-badge">关键词匹配</span>' : '';
                
                fileItem.innerHTML = `
                    <i class="${fileIcon} file-icon"></i>
                    <div class="file-info">
                        <div class="file-name">${file.name}</div>
                        <div class="file-meta">
                            <span>大小: ${fileSize}</span>
                            <span>•</span>
                            <span>路径: ${file.path}</span>
                        </div>
                    </div>
                    ${badge}
                    <div class="file-actions">
                        <button class="file-action-btn" onclick="openFile('${encodeURIComponent(file.path)}', '${file.name}')" title="打开文件">
                            <i class="fas fa-external-link-alt"></i>
                        </button>
                    </div>
                `;
                
                // 添加点击事件来打开文件
                fileItem.addEventListener('click', (e) => {
                    // 如果点击的是按钮，不触发文件打开
                    if (e.target.closest('.file-action-btn')) {
                        return;
                    }
                    openFile(encodeURIComponent(file.path), file.name);
                });
                
                fileList.appendChild(fileItem);
            });
        } else {
            fileList.innerHTML = '<p style="text-align: center; color: #666;">暂无文件</p>';
        }
        
        fileListModal.style.display = 'flex';
        fileListModal.classList.add('show');
        modalOverlay.style.display = 'block';
        modalOverlay.classList.add('show');
        
    } catch (error) {
        console.error('获取文件列表失败:', error);
        alert('获取文件列表失败');
    }
}

// 打开文件函数
async function openFile(filePath, fileName) {
    try {
        // 构建文件访问URL
        const fileUrl = `/files/${filePath}`;
        
        // 获取文件扩展名
        const fileExt = fileName.split('.').pop().toLowerCase();
        
        // 定义可以在浏览器中直接打开的文件类型
        const browserDisplayableTypes = [
            'pdf', 'jpg', 'jpeg', 'png', 'gif', 'bmp', 'svg',
            'txt', 'html', 'css', 'js', 'json', 'xml',
            'mp4', 'avi', 'mov', 'wmv', 'flv',
            'mp3', 'wav', 'ogg', 'aac',
            'csv', 'md', 'log'
            // 移除Office文档，因为浏览器通常不支持直接显示
            // 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx'
        ];
        
        // 定义Office文档类型
        const officeTypes = ['doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx'];
        
        // 检查文件是否可以在浏览器中直接显示
        const canDisplayInBrowser = browserDisplayableTypes.includes(fileExt);
        const isOfficeDocument = officeTypes.includes(fileExt);
        
        if (canDisplayInBrowser) {
            // 可以在浏览器中显示的文件，直接在新窗口中打开
            console.log('在浏览器中打开文件:', fileName);
            window.open(fileUrl, '_blank');
        } else if (isOfficeDocument) {
            // Office文档，提供在线预览选项
            const choice = confirm(`${fileName} 是Office文档。\n\n选择操作：\n- 确定：尝试在线预览\n- 取消：直接下载`);
            
            if (choice) {
                // 尝试使用在线预览服务
                const previewUrl = `https://view.officeapps.live.com/op/embed.aspx?src=${encodeURIComponent(window.location.origin + fileUrl)}`;
                console.log('尝试在线预览Office文档:', fileName);
                window.open(previewUrl, '_blank');
            } else {
                // 直接下载
                console.log('下载Office文档:', fileName);
                const downloadUrl = `/download/${filePath}`;
                const link = document.createElement('a');
                link.href = downloadUrl;
                link.download = fileName;
                link.target = '_blank';
                document.body.appendChild(link);
                link.click();
                document.body.removeChild(link);
            }
        } else {
            // 不能在浏览器中显示的文件，提供下载
            console.log('下载文件:', fileName);
            const downloadUrl = `/download/${filePath}`;
            const link = document.createElement('a');
            link.href = downloadUrl;
            link.download = fileName;
            link.target = '_blank';
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
        }
    } catch (error) {
        console.error('打开文件失败:', error);
        alert(`无法打开文件: ${fileName}`);
    }
}

// 关闭模态框
function closeModal() {
    fileListModal.style.display = 'none';
    fileListModal.classList.remove('show');
    modalOverlay.style.display = 'none';
    modalOverlay.classList.remove('show');
}

// 获取文件图标
function getFileIcon(filename) {
    const ext = filename.split('.').pop().toLowerCase();
    
    const iconMap = {
        'pdf': 'fas fa-file-pdf',
        'doc': 'fas fa-file-word',
        'docx': 'fas fa-file-word',
        'xls': 'fas fa-file-excel',
        'xlsx': 'fas fa-file-excel',
        'ppt': 'fas fa-file-powerpoint',
        'pptx': 'fas fa-file-powerpoint',
        'txt': 'fas fa-file-alt',
        'jpg': 'fas fa-file-image',
        'jpeg': 'fas fa-file-image',
        'png': 'fas fa-file-image',
        'gif': 'fas fa-file-image',
        'zip': 'fas fa-file-archive',
        'rar': 'fas fa-file-archive',
        'mp4': 'fas fa-file-video',
        'avi': 'fas fa-file-video',
        'mp3': 'fas fa-file-audio',
        'wav': 'fas fa-file-audio'
    };
    
    return iconMap[ext] || 'fas fa-file';
}

// 格式化文件大小
function formatFileSize(bytes) {
    if (bytes === 0) return '0 B';
    
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}

// 处理 ESC 键关闭模态框
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && fileListModal.style.display === 'flex') {
        closeModal();
    }
});