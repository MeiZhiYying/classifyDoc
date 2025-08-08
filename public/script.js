// 全局变量
let currentStats = {};
let allFiles = [];
let currentSort = 'time';
let currentOrder = 'desc';
let currentFilter = '';

// 分类配置
const CATEGORY_CONFIG = {
    '合同': {
        icon: 'fas fa-file-contract',
        color: 'contract',
        description: '包含合同协议等相关文件'
    },
    '简历': {
        icon: 'fas fa-user-tie',
        color: 'resume',
        description: '包含个人简历和求职相关文件'
    },
    '发票': {
        icon: 'fas fa-receipt',
        color: 'invoice',
        description: '包含发票收据等财务文件'
    },
    '论文': {
        icon: 'fas fa-graduation-cap',
        color: 'thesis',
        description: '包含学术论文和研究文件'
    },
    '未分类': {
        icon: 'fas fa-question-circle',
        color: 'unclassified',
        description: '暂未识别类型的文件'
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

// 文件列表相关元素
const categoryFilter = document.getElementById('categoryFilter');
const sortByTime = document.getElementById('sortByTime');
const sortBySize = document.getElementById('sortBySize');
const fileListTable = document.getElementById('fileListTable');
const fileListBody = document.getElementById('fileListBody');
const noFilesMessage = document.getElementById('noFilesMessage');

// 新增分类模态框元素
const addCategoryModal = document.getElementById('addCategoryModal');
const categoryNameInput = document.getElementById('categoryName');
const usernameInput = document.getElementById('username');

// 初始化
document.addEventListener('DOMContentLoaded', function() {
    initializeUpload();
    initializeFileList();
    initializeTheme();
    initializeAuthUI();
    initializeCommonSettings();
    registerGlobalDismiss();
    // 先渲染一次默认分类，保证首屏就可见
    try { renderCategories(); } catch (e) { console.warn('首屏分类渲染失败(忽略):', e); }
    if (shouldAutoScanUploads()) {
        autoScanUploads();
    }
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
        } else {
            console.error('扫描失败:', result.error);
        }
        
        // 加载统计数据和文件列表
        loadStats();
        loadAllFiles();
        
    } catch (error) {
        console.error('自动扫描失败:', error);
        // 即使扫描失败，也要加载统计数据和文件列表
        loadStats();
        loadAllFiles();
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
            loadAllFiles(); // 刷新文件列表
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

// 更新文件数量显示
function updateFileCount(count) {
    const fileCount = document.getElementById('fileCount');
    if (fileCount) {
        fileCount.textContent = `(${count}个)`;
    }
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
        renderCategories();
    }
}

// 渲染分类卡片 - 固定8个位置的2x4布局
function renderCategories() {
    console.log('开始渲染分类卡片，currentStats:', currentStats);
    categoriesGrid.innerHTML = '';
    
    // 确保currentStats不为空，如果为空则使用默认值
    if (!currentStats || Object.keys(currentStats).length === 0) {
        currentStats = {
            '合同': { count: 0, files: [] },
            '简历': { count: 0, files: [] },
            '发票': { count: 0, files: [] },
            '论文': { count: 0, files: [] },
            '未分类': { count: 0, files: [] }
        };
    }
    
    // 固定的分类顺序：前4个预定义分类
    const fixedCategories = ['合同', '简历', '发票', '论文'];
    
    // 获取新增的分类（排除预定义分类和未分类）
    const allCategories = Object.keys(currentStats);
    const newCategories = allCategories.filter(cat => 
        !fixedCategories.includes(cat) && cat !== '未分类'
    ).slice(0, 3); // 最多3个新增分类
    
    // 创建固定的8个位置
    const positions = [
        // 第一行：4个预定义分类
        ...fixedCategories,
        // 第二行：3个新增分类位置 + 1个未分类
        ...Array(3).fill(null).map((_, index) => 
            newCategories[index] || `add-new-${index}`
        ),
        '未分类'
    ];
    
    positions.forEach((position, index) => {
        if (position.startsWith('add-new-')) {
            // 创建"新增分类"卡片
            const addNewIndex = parseInt(position.split('-')[2]);
            const card = document.createElement('div');
            card.className = 'category-card add-new';
            card.innerHTML = `
                <i class="fas fa-plus category-icon"></i>
                <h3 class="category-title">新增分类</h3>
                <div class="category-count">+</div>
                <p class="category-description">点击添加新分类</p>
            `;
            
            // 添加点击事件
            card.addEventListener('click', () => {
                showAddCategoryModal();
            });
            
            categoriesGrid.appendChild(card);
        } else {
            // 创建正常分类卡片
            const categoryName = position;
            const stats = currentStats[categoryName] || { count: 0, files: [] };
            
            // 获取分类配置
            let config = CATEGORY_CONFIG[categoryName];
            if (!config) {
                // 为新分类创建默认配置
                config = {
                    icon: 'fas fa-folder',
                    color: 'default',
                    description: '自定义分类'
                };
            }
            
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
        }
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
                        <button class="file-action-btn" onclick="openFile('${encodeURIComponent(file.path)}', '${file.name}')" title="下载文件">
                            <i class="fas fa-download"></i>
                        </button>
                    </div>
                `;
                
                // 添加点击事件来下载文件
                fileItem.addEventListener('click', (e) => {
                    // 如果点击的是按钮，不触发文件下载
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

// 下载文件函数
async function openFile(filePath, fileName) {
    try {
        // 弹出确认对话框，点击确定下载
        const confirmDownload = confirm(`是否确定下载文件: ${fileName}？`);
        
        if (confirmDownload) {
            // 用户点击确定，执行下载
            console.log('下载文件:', fileName);
            const downloadUrl = `/download/${filePath}`;
            const link = document.createElement('a');
            link.href = downloadUrl;
            link.download = fileName;
            link.target = '_blank';
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
        } else {
            // 用户点击取消，不执行任何操作
            console.log('用户取消下载:', fileName);
        }
    } catch (error) {
        console.error('下载文件失败:', error);
        alert(`无法下载文件: ${fileName}`);
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

// 显示新增分类模态框
function showAddCategoryModal() {
    document.getElementById('addCategoryModal').style.display = 'flex';
}

// 关闭新增分类模态框
function closeAddCategoryModal() {
    document.getElementById('addCategoryModal').style.display = 'none';
    // 清空表单
    document.getElementById('categoryName').value = '';
    document.getElementById('username').value = '';
}

// 添加新分类
async function addCategory() {
    const categoryName = document.getElementById('categoryName').value.trim();
    const username = document.getElementById('username').value.trim();
    
    if (!categoryName || !username) {
        alert('请填写完整的分类名称和用户名');
        return;
    }
    
    try {
        console.log('正在添加新分类:', { categoryName, username });
        
        const response = await fetch('/api/add-category', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                categoryName: categoryName,
                username: username
            })
        });
        
        console.log('服务器响应状态:', response.status);
        
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}`);
        }
        
        const result = await response.json();
        console.log('服务器响应:', result);
        
        if (result.success) {
            alert(result.message);
            closeAddCategoryModal();
            
            // 等待一段时间让后台扫描完成
            console.log('等待后台扫描完成...');
            setTimeout(() => {
                // 重新加载统计数据以显示新分类
                loadStats();
                loadAllFiles();
                console.log('重新加载完成');
            }, 2000);
        } else {
            alert('添加分类失败: ' + result.error);
            console.error('添加分类失败:', result.error);
        }
        
    } catch (error) {
        console.error('添加分类失败:', error);
        alert('添加分类失败: ' + error.message);
    }
}

// 处理 ESC 键关闭模态框
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        if (fileListModal.style.display === 'flex') {
            closeModal();
        }
        if (document.getElementById('addCategoryModal').style.display === 'flex') {
            closeAddCategoryModal();
        }
        const settingsDropdown = document.getElementById('settingsDropdown');
        if (settingsDropdown && !settingsDropdown.hidden) {
            hideSettings();
        }
        const accountDropdown = document.getElementById('accountDropdown');
        if (accountDropdown && !accountDropdown.hidden) {
            hideAccount();
        }
        const authModal = document.getElementById('authModal');
        if (authModal && authModal.style.display === 'flex') {
            closeAuthModal();
        }
    }
});

// === 主题设置 ===
function initializeTheme() {
    const settingsButton = document.getElementById('settingsButton');
    const settingsDropdown = document.getElementById('settingsDropdown');
    const themeLight = document.getElementById('themeLight');
    const themeDark = document.getElementById('themeDark');

    if (!settingsButton || !settingsDropdown || !themeLight || !themeDark) {
        return;
    }

    // 加载主题
    const savedTheme = getSavedTheme();
    applyTheme(savedTheme || 'light');
    if (savedTheme === 'dark') {
        themeDark.checked = true;
    } else {
        themeLight.checked = true;
    }

    // 切换设置面板（回退到简单显示/隐藏，不改动布局）
    settingsButton.addEventListener('click', (e) => {
        e.stopPropagation();
        // 打开设置前，确保关闭账户下拉，避免重叠
        hideAccount();
        const expanded = settingsButton.getAttribute('aria-expanded') === 'true';
        const dropdown = document.getElementById('settingsDropdown');
        if (!dropdown) return;
        dropdown.hidden = expanded;
        settingsButton.setAttribute('aria-expanded', expanded ? 'false' : 'true');
    });

    // 点击外部关闭
    document.addEventListener('click', (e) => {
        if (!settingsDropdown.hidden && !settingsDropdown.contains(e.target) && e.target !== settingsButton) {
            hideSettings();
        }
    });

    // 主题选择
    themeLight.addEventListener('change', () => setTheme('light'));
    themeDark.addEventListener('change', () => setTheme('dark'));
}

function setTheme(theme) {
    applyTheme(theme);
    saveTheme(theme);
}

function applyTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
}

function saveTheme(theme) {
    try {
        localStorage.setItem('app_theme', theme);
    } catch (e) {}
}

function getSavedTheme() {
    try {
        return localStorage.getItem('app_theme');
    } catch (e) { return null; }
}

function showSettings() {
    const dropdown = document.getElementById('settingsDropdown');
    const btn = document.getElementById('settingsButton');
    if (dropdown && btn) {
        dropdown.hidden = false;
        btn.setAttribute('aria-expanded', 'true');
    }
}

function hideSettings() {
    const dropdown = document.getElementById('settingsDropdown');
    const btn = document.getElementById('settingsButton');
    if (dropdown && btn) {
        dropdown.hidden = true;
        btn.setAttribute('aria-expanded', 'false');
    }
}

// === 账户/鉴权 ===
function initializeAuthUI() {
    const accountButton = document.getElementById('accountButton');
    const accountDropdown = document.getElementById('accountDropdown');
    const loginOpenBtn = document.getElementById('loginOpenBtn');
    const logoutBtn = document.getElementById('logoutBtn');

    if (!accountButton || !accountDropdown) return;

    // 首次查询登录状态
    refreshAuthState();

    accountButton.addEventListener('click', (e) => {
        e.stopPropagation();
        // 打开账户前，确保关闭设置下拉，避免重叠
        hideSettings();
        const expanded = accountButton.getAttribute('aria-expanded') === 'true';
        const dropdown = document.getElementById('accountDropdown');
        if (!dropdown) return;
        dropdown.hidden = expanded;
        accountButton.setAttribute('aria-expanded', expanded ? 'false' : 'true');
    });

    document.addEventListener('click', (e) => {
        if (!accountDropdown.hidden && !accountDropdown.contains(e.target) && e.target !== accountButton) {
            hideAccount();
        }
    });

    if (loginOpenBtn) {
        loginOpenBtn.addEventListener('click', () => {
            hideAccount();
            openAuthModal();
        });
    }
    if (logoutBtn) {
        logoutBtn.addEventListener('click', async () => {
            await fetch('/api/auth/logout', { method: 'POST' });
            refreshAuthState();
        });
    }

    // 登录模态事件
    const authModal = document.getElementById('authModal');
    const authCloseBtn = document.getElementById('authCloseBtn');
    const authCancelBtn = document.getElementById('authCancelBtn');
    const authSubmitBtn = document.getElementById('authSubmitBtn');
    authCloseBtn && authCloseBtn.addEventListener('click', closeAuthModal);
    authCancelBtn && authCancelBtn.addEventListener('click', closeAuthModal);
    authSubmitBtn && authSubmitBtn.addEventListener('click', submitLogin);

    // 在登录模态内支持“注册/登录”切换：按回车登录
    const loginPassword = document.getElementById('loginPassword');
    loginPassword && loginPassword.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') submitLogin();
    });

    function submitLogin() {
        const username = document.getElementById('loginUsername').value.trim();
        const password = document.getElementById('loginPassword').value;
        if (!username || !password) {
            alert('请输入用户名和密码');
            return;
        }
        const doLogin = () => fetch('/api/auth/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        }).then(async res => {
            if (!res.ok) {
                const t = await res.json().catch(() => ({}));
                throw new Error(t.error || '登录失败');
            }
            return res.json();
        }).then(() => {
            closeAuthModal();
            refreshAuthState();
        });

        // 先尝试注册（若用户不存在则注册），再登录
        fetch('/api/auth/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        }).finally(() => {
            doLogin().catch(err => alert(err.message || '登录失败'))
        });
    }

    function openAuthModal() {
        if (authModal) {
            authModal.style.display = 'flex';
            authModal.classList.add('show');
        }
    }
    function closeAuthModal() {
        if (authModal) {
            authModal.style.display = 'none';
            authModal.classList.remove('show');
        }
    }
    // 暴露给 Esc 关闭
    window.closeAuthModal = closeAuthModal;
    window.openAuthModal = openAuthModal;
}

function showAccount() {
    const dropdown = document.getElementById('accountDropdown');
    const btn = document.getElementById('accountButton');
    if (dropdown && btn) {
        dropdown.hidden = false;
        btn.setAttribute('aria-expanded', 'true');
    }
}
function hideAccount() {
    const dropdown = document.getElementById('accountDropdown');
    const btn = document.getElementById('accountButton');
    if (dropdown && btn) {
        dropdown.hidden = true;
        btn.setAttribute('aria-expanded', 'false');
    }
}

async function refreshAuthState() {
    try {
        const res = await fetch('/api/auth/me');
        const data = await res.json();
        const info = document.getElementById('accountInfo');
        const loginBtn = document.getElementById('loginOpenBtn');
        const logoutBtn = document.getElementById('logoutBtn');
        if (data.authenticated) {
            info.textContent = `已登录：${data.user.username}`;
            loginBtn.hidden = true;
            logoutBtn.hidden = false;
        } else {
            info.textContent = '未登录';
            loginBtn.hidden = false;
            logoutBtn.hidden = true;
        }
    } catch (e) {
        // 忽略错误
    }
}

// 常用设置：是否启动时自动扫描 uploads
function initializeCommonSettings() {
    const toggle = document.getElementById('autoScanToggle');
    if (!toggle) return;
    const key = 'setting_auto_scan_uploads';
    try {
        const saved = localStorage.getItem(key);
        if (saved !== null) {
            toggle.checked = saved === '1';
        } else {
            toggle.checked = true; // 默认开启
            localStorage.setItem(key, '1');
        }
    } catch (e) {}

    toggle.addEventListener('change', () => {
        try {
            localStorage.setItem(key, toggle.checked ? '1' : '0');
        } catch (e) {}
    });
}

function shouldAutoScanUploads() {
    try {
        const v = localStorage.getItem('setting_auto_scan_uploads');
        return v === null || v === '1';
    } catch (e) { return true; }
}

// 全局点击空白关闭所有下拉
function registerGlobalDismiss() {
    document.addEventListener('click', (e) => {
        const settingsDropdown = document.getElementById('settingsDropdown');
        const settingsButton = document.getElementById('settingsButton');
        const accountDropdown = document.getElementById('accountDropdown');
        const accountButton = document.getElementById('accountButton');

        const clickInSettings = settingsDropdown && (settingsDropdown.contains(e.target) || (settingsButton && settingsButton.contains(e.target)));
        const clickInAccount = accountDropdown && (accountDropdown.contains(e.target) || (accountButton && accountButton.contains(e.target)));

        if (!clickInSettings) {
            hideSettings();
        }
        if (!clickInAccount) {
            hideAccount();
        }
    });
}

// === 文件列表功能 ===

// 初始化文件列表功能
function initializeFileList() {
    // 初始化分类筛选器
    updateCategoryFilter();
    
    // 绑定筛选器事件
    categoryFilter.addEventListener('change', handleCategoryFilter);
    
    // 绑定排序按钮事件
    sortByTime.addEventListener('click', () => handleSort('time'));
    sortBySize.addEventListener('click', () => handleSort('size'));
    
    // 加载文件列表
    loadAllFiles();
}

// 更新分类筛选器选项
function updateCategoryFilter() {
    // 清空现有选项
    categoryFilter.innerHTML = '<option value="">全部类型</option>';
    
    // 添加分类选项
    Object.keys(CATEGORY_CONFIG).forEach(categoryName => {
        const option = document.createElement('option');
        option.value = categoryName;
        option.textContent = categoryName;
        categoryFilter.appendChild(option);
    });
}

// 处理分类筛选
function handleCategoryFilter() {
    currentFilter = categoryFilter.value;
    loadAllFiles();
}

// 处理排序
function handleSort(sortType) {
    // 如果点击的是当前激活的排序，则切换排序顺序
    if (currentSort === sortType) {
        currentOrder = currentOrder === 'desc' ? 'asc' : 'desc';
    } else {
        // 如果切换到新的排序方式，默认使用降序
        currentSort = sortType;
        currentOrder = sortType === 'time' ? 'desc' : 'desc'; // 时间默认降序（最新在前），大小默认降序（大的在前）
    }
    
    // 更新按钮状态
    updateSortButtons();
    
    // 重新加载文件列表
    loadAllFiles();
}

// 更新排序按钮状态
function updateSortButtons() {
    // 重置所有按钮
    sortByTime.classList.remove('active', 'desc', 'asc');
    sortBySize.classList.remove('active', 'desc', 'asc');
    
    // 更新图标
    sortByTime.querySelector('.sort-icon').className = 'fas fa-sort sort-icon';
    sortBySize.querySelector('.sort-icon').className = 'fas fa-sort sort-icon';
    
    // 设置激活的按钮
    if (currentSort === 'time') {
        sortByTime.classList.add('active', currentOrder);
        sortByTime.querySelector('.sort-icon').className = currentOrder === 'desc' ? 'fas fa-sort-down sort-icon' : 'fas fa-sort-up sort-icon';
    } else if (currentSort === 'size') {
        sortBySize.classList.add('active', currentOrder);
        sortBySize.querySelector('.sort-icon').className = currentOrder === 'desc' ? 'fas fa-sort-down sort-icon' : 'fas fa-sort-up sort-icon';
    }
}

// 加载所有文件列表
async function loadAllFiles() {
    try {
        const params = new URLSearchParams({
            sort: currentSort,
            order: currentOrder
        });
        
        if (currentFilter) {
            params.append('category', currentFilter);
        }
        
        const response = await fetch(`/api/all-files?${params}`);
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}`);
        }
        
        const data = await response.json();
        allFiles = data.files || [];
        
        renderFileList();
        
    } catch (error) {
        console.error('加载文件列表失败:', error);
        allFiles = [];
        renderFileList();
    }
}

// 渲染文件列表
function renderFileList() {
    // 清空表格内容
    fileListBody.innerHTML = '';
    
    if (allFiles.length === 0) {
        // 显示无文件消息
        fileListTable.style.display = 'none';
        noFilesMessage.style.display = 'block';
        return;
    }
    
    // 隐藏无文件消息，显示表格
    noFilesMessage.style.display = 'none';
    fileListTable.style.display = 'table';
    
    // 更新文件数量显示
    updateFileCount(allFiles.length);
    
    // 渲染文件行
    allFiles.forEach(file => {
        const row = createFileRow(file);
        fileListBody.appendChild(row);
    });
}

// 创建文件行
function createFileRow(file) {
    const row = document.createElement('tr');
    
    // 文件图标
    const iconCell = document.createElement('td');
    iconCell.innerHTML = `<i class="${getFileIcon(file.name)} file-table-icon"></i>`;
    row.appendChild(iconCell);
    
    // 文件名
    const nameCell = document.createElement('td');
    nameCell.innerHTML = `<div class="file-table-name">${file.name}</div>`;
    row.appendChild(nameCell);
    
    // 分类
    const categoryCell = document.createElement('td');
    const categoryConfig = CATEGORY_CONFIG[file.category] || CATEGORY_CONFIG['未分类'];
    categoryCell.innerHTML = `<span class="file-table-category ${categoryConfig.color}">${file.category}</span>`;
    row.appendChild(categoryCell);
    
    // 大小
    const sizeCell = document.createElement('td');
    sizeCell.innerHTML = `<div class="file-table-size">${formatFileSize(file.size)}</div>`;
    row.appendChild(sizeCell);
    
    // 时间
    const timeCell = document.createElement('td');
    const timeStr = formatFileTime(file.modTime);
    timeCell.innerHTML = `<div class="file-table-time">${timeStr}</div>`;
    row.appendChild(timeCell);
    
    // 操作
    const actionCell = document.createElement('td');
    actionCell.innerHTML = `<button class="file-table-action" onclick="openFile('${encodeURIComponent(file.path)}', '${file.name}')" title="下载文件">
        <i class="fas fa-download"></i>
    </button>`;
    row.appendChild(actionCell);
    
    return row;
}

// 格式化文件时间
function formatFileTime(modTime) {
    const date = new Date(modTime);
    const now = new Date();
    const diffMs = now - date;
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
    
    if (diffHours < 1) {
        return '刚刚';
    } else if (diffHours < 24) {
        return `${diffHours}小时前`;
    } else if (diffDays < 7) {
        return `${diffDays}天前`;
    } else {
        return date.toLocaleDateString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit'
        });
    }
}