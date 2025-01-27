document.addEventListener('DOMContentLoaded', function () {
  const stepsContainer = document.querySelector('.steps');
  //const addStepBtn = document.getElementById('add-step-btn');
  const dropdownContent = document.getElementById('dropdown-content');
  const stepInfoDiv = document.getElementById('step-info'); // 获取步骤信息显示区域
  let steps = [];  // 保存步骤信息的数组

  function genuuid() {
    if (typeof crypto !== 'undefined' && crypto.randomUUID) {
      return crypto.randomUUID();
    } else {
      // 处理不支持 crypto.randomUUID() 的环境
      console.warn("crypto.randomUUID() is not supported in this environment, falling back to less secure method.");
      return generateUUIDv4Fallback();
    }
  }

  // 创建步骤的函数
  function createStep(stepType = "default") {
    const step = document.createElement('div');
    step.classList.add('step');
    step.dataset.stepid = genuuid();
    step.dataset.type = stepType;
    step.dataset.disabled = false;

    switch (stepType) {
      case "filter-include":
        step.innerHTML = `
        <div class="step-header">
          过滤-白名单
          <div>
            <label class="switch">
              <input type="checkbox" class="switch-checkbox" checked>
              <span class="slider"></span>
            </label>
            <button class="remove-step-btn">删除</button>
          </div>
        </div>
        <div class="step-content">
          <div class="inline-container">
            <select class="mode-sel">
              <option value="equal">等于</option>
              <option value="contain">包含</option>
              <option value="prefix">前缀</option>
              <option value="suffix">后缀</option>
              <option value="regexp">正则</option>
            </select>
            <input type="text" class="auto" placeholder="请输入表达式">
            <input type="text" class="mini" placeholder="过滤次数">
            <button class="preview-btn">预览</button>
          </div>
          <button class="prev-btn">上一步</button>
          <button class="next-btn">下一步</button>
          <button class="submit-btn">提交</button>
        </div>
        <div class="step-result"></div>
        <div class="step-preview">
          <div class="preview-title">
            <h3>效果预览</h3>
            <div class="step-result"></div>
          </div>
          <div class="preview-container">
          </div>
        </div>`;
        break;
      case "filter-exclude":
        step.innerHTML = `
        <div class="step-header">
          过滤-黑名单
          <div>
            <label class="switch">
              <input type="checkbox" class="switch-checkbox" checked>
              <span class="slider"></span>
            </label>
            <button class="remove-step-btn">删除</button>
          </div>
        </div>
        <div class="step-content">
          <div class="inline-container">
            <select class="mode-sel">
              <option value="equal">等于</option>
              <option value="contain">包含</option>
              <option value="prefix">前缀</option>
              <option value="suffix">后缀</option>
              <option value="regexp">正则</option>
            </select>
            <input type="text" class="auto" placeholder="请输入表达式">
            <input type="text" class="mini" placeholder="过滤次数">
            <button class="preview-btn">预览</button>
          </div>
          <button class="prev-btn">上一步</button>
          <button class="next-btn">下一步</button>
          <button class="submit-btn">提交</button>
        </div>
        <div class="step-preview">
          <div class="preview-title">
            <h3>效果预览</h3>
            <div class="step-result"></div>
          </div>
          <div class="preview-container">
          </div>
        </div>`;
        break;
      case "rename-add":
        step.innerHTML = `
        <div class="step-header">
          重命名-添加
          <div>
            <label class="switch">
              <input type="checkbox" class="switch-checkbox" checked>
              <span class="slider"></span>
            </label>
            <button class="remove-step-btn">删除</button>
          </div>
          </div>
          <div class="step-content">
            <div class="inline-container">
              <select class="mode-sel">
                <option value="plain">文本</option>
                <option value="prefix">前缀</option>
                <option value="suffix">后缀</option>
                <option value="index">索引</option>
                <option value="regexp">正则</option>
              </select>
              <input type="text" class="auto" placeholder="请输入表达式">
              <input type="text" class="auto" placeholder="请输入添加值">
              <input type="text" class="mini" placeholder="匹配次数">
              <button class="preview-btn">预览</button>
            </div>
            <button class="prev-btn">上一步</button>
            <button class="next-btn">下一步</button>
            <button class="submit-btn">提交</button>
          </div>
          <div class="step-preview">
            <div class="preview-title">
              <h3>效果预览</h3>
              <div class="step-result"></div>
            </div>
            <div class="preview-container">
            </div>
          </div>`;
        break;
      default:
        if (steps.length > 0) { // 非第一个步骤
          return;
        }
        // step.classList.add('active');
        // selectedValue.dataset.fid = '0';
        // selectedValue.dataset.tooltip = '/';
        // selectedValue.textContent = '根目录';
        step.innerHTML = `
        <div class="step-header">
            文件选择
        </div>
        <div class="step-content">
          <div class="inline-container" data-type="default">
            <div class="tree-select">
              <div class="selected-value" data-fid="0" data-tooltip="/" >根目录</div>
              <ul class="tree-options">
                <!-- 初始数据将在这里动态生成 -->
              </ul>
            </div>
            <input type="text" class="search-value auto" placeholder="请输入搜索关键字(可选)">
            <button class="preview-btn">预览</button>
          </div>
          <button class="prev-btn">上一步</button>
          <button class="next-btn">下一步</button>
        </div>
        
        <div class="step-preview">
          <div class="preview-title">
            <h3>效果预览</h3>
            <div class="step-result"></div>
          </div>
          <div class="preview-container">
          </div>
        </div>`;
        break;
    }

    return step;
  }

  function renderSteps(stepid) {
    stepsContainer.innerHTML = ''; // 清空步骤容器
    steps.forEach(step => {
      stepsContainer.appendChild(step);
    });

    showStep(stepid);

    updateStepInfo(); // 更新步骤信息
  }

  function addStep(stepType) {
    const step = createStep(stepType);
    if (step === undefined) {
      return;
    }

    let index = stepIndex();
    // 当前步骤后添加步骤
    steps.splice(index + 1, 0, step);
    // 渲染步骤
    renderSteps(step.dataset.stepid);
    // 添加步骤事件
    setupStepEvents(step);
  }

  function removeStep(stepid) {
    let idx = stepIndex(stepid);
    if (idx <= 0) {
      return;
    }

    steps.splice(idx, 1);
    renderSteps(steps[idx - 1].dataset.stepid);
  }

  // 刷新preview组件
  function refreshPreview(preview, data) {
    const previewContainer = preview.querySelector('.preview-container');
    previewContainer.innerHTML = '';
    if (data === null || data.length === 0) {
      return;
    }

    data.forEach(file => {
      if (!preview.closest('.step').dataset.type.startsWith("rename-")) {
        previewContainer.innerHTML += `
        <div class="file-item" data-fid="${file.fileId}" data-ftype="${file.type}" data-fpath="${file.path}" data-fname="${file.filename}" data-nfname="${file.nfname}" data-selected="1">
          <div class="file-type-indicator"></div>
          <span class="file-path">${file.path.substring(0, file.path.lastIndexOf('/') + 1)}</span>
          <div class="file-name-container">
            <div class="file-name-row">
              <span class="file-name-label original"><span>原：</span></span>
              <span class="file-name original">${file.nfname}</span>
            </div>
          </div>
        </div>
        `;
      } else {
        let ostr = "";
        let nstr = "";

        var dmp = new diff_match_patch();
        var diff = dmp.diff_main(file.filename, file.nfname);
        diff.forEach((item) => {
          // -1：原始地址
          if (item[0] === -1) {
            ostr += `<span style = 'background:#f6c5c5'>${item[1]}</span>`;
          }
          // 1：匹配地址
          else if (item[0] === 1) {
            nstr += `<span style = 'background:#b5ecb5'>${item[1]}</span>`;
          } else {
            //1：共有
            ostr += `<span>${item[1]}</span>`;
            nstr += `<span>${item[1]}</span>`;
          }
        });

        previewContainer.innerHTML += `
        <div class="file-item" data-fid="${file.fileId}" data-ftype="${file.type}" data-fpath="${file.path}" data-fname="${file.filename}" data-nfname="${file.nfname}" data-selected="1">
          <div class="file-type-indicator"></div>
          <span class="file-path">${file.path.substring(0, file.path.lastIndexOf('/') + 1)}</span>
          <div class="file-name-container">
            <div class="file-name-row">
              <span class="file-name-label original"><span>原：</span></span>
              <span class="file-name original">${ostr}</span>
            </div>
            <div class="file-name-row">
              <span class="file-name-label"><span>新：</span></span>
              <span class="file-name">${nstr}</span>
            </div>
          </div>
        </div>
        `;
      }
    });

    const fitems = preview.querySelectorAll('.file-item');
    fitems.forEach(item => {
      item.addEventListener('click', function (event) {
        event.stopPropagation(); // 阻止事件冒泡到父级 li 和 <treeSelect></treeSelect>
        if (item.dataset.selected == "1") {
          item.dataset.selected = "0";
        } else {
          item.dataset.selected = "1";
        }
      });

      item.addEventListener('dblclick', function (event) {
        event.stopPropagation(); // 阻止事件冒泡到父级 li 和 <treeSelect></treeSelect>
        const istep = this.closest('.step');
        if (istep.dataset.type === "default") {
          // 更新selectedValue
          const selectedValue = istep.querySelector('.tree-select').querySelector('.selected-value');
          selectedValue.dataset.fid = this.dataset.fid;
          selectedValue.dataset.tooltip = this.dataset.fpath;
          selectedValue.textContent = this.dataset.fpath;
        } else if (istep.dataset.type.startsWith("filter-")) {
          // 更新selectedValue
          const select = istep.querySelector('.inline-container').querySelector('.mode-sel');
          select.value = "equal";
          const input = istep.querySelector('.inline-container').querySelector('.auto');
          input.value = this.dataset.fname;
          const viewbtn = istep.querySelector('.inline-container').querySelector('.preview-btn');
          viewbtn.click();
        }
        // console.log(this.dataset.fid, this.dataset.fpath, this.dataset.fname, this.dataset.nfname);
      });
    });
  }

  // 继承上一步骤数据更新preview
  function updatePreviewByInherit(preview) {
    let prestep = preview.closest('.step');
    let idx = stepIndex(prestep.dataset.stepid);
    if (idx <= 0) {
      return;
    }
    prestep = steps[idx - 1];
    while (prestep.dataset.disabled === "true") {
      idx--;
      prestep = steps[idx - 1];
    }

    const previewContainer = prestep.querySelector('.preview-container');
    const fileitems = previewContainer.querySelectorAll('.file-item');
    let data = [];
    fileitems.forEach(item => {
      if (item.dataset.selected == "1") {
        data.push({
          fileId: item.dataset.fid,
          filename: item.dataset.nfname,
          type: item.dataset.ftype,
          path: item.dataset.fpath,
          nfname: item.dataset.nfname
        });
      }
    });

    // 更新data
    const cstep = preview.closest('.step');
    let flag = false;
    cstep.querySelectorAll('input.auto').forEach(input => {
      if (input.value != "") {
        flag = true;
      }
    });

    if (flag) {
      let kind = cstep.dataset.type;
      kind = kind.substring(kind.lastIndexOf('-') + 1);
      let mode = cstep.querySelector('select').value;
      let inputs = cstep.querySelector(".inline-container").querySelectorAll('input');

      let expr = "";
      let value = "";
      let num = "";
      if (inputs.length == 1) {
        expr = inputs[0].value;
      }
      if (inputs.length == 2) {
        expr = inputs[0].value;
        num = inputs[1].value;
      }
      if (inputs.length == 3) {
        expr = inputs[0].value;
        value = inputs[1].value;
        num = inputs[2].value;
      }

      let action = new Action(kind, mode, expr, value, num);
      if (kind == "include" || kind == "exclude") {
        strs = action.filter(data.map(item => item.nfname));
        data = data.filter(item => strs.includes(item.nfname));
      } else {
        strs = action.rename(data.map(item => item.nfname));
        data = data.map((item, index) => {
          item.filename = item.nfname;
          item.nfname = strs[index];
          return item;
        });
      }
    }

    refreshPreview(preview, data);
  }

  // 请求后端数据更新preview
  function updatePreviewByFetch(preview, cid, keyword) {
    return fetch('/list/p123', { method: 'POST', body: JSON.stringify({ cid: cid.toString(), key: keyword.toString() }) }).then(response => response.json()).then(data => {
      if (data != null) {
        data = data.map((item) => {
          item.nfname = item.filename;
          return item;
        });
      }
      refreshPreview(preview, data);
    }).catch(error => {
      console.error('Error:', error);
    });
  }

  // 初始化加载文件选择下拉树
  async function initializeTree() {
    const treeSelect = document.querySelector('.tree-select');
    const selectedValue = treeSelect.querySelector('.selected-value');
    const searchValue = treeSelect.closest('.inline-container').querySelector('.search-value');
    // 创建 MutationObserver 实例
    const observer = new MutationObserver(mutations => {
      mutations.forEach(mutation => {
        if (mutation.type === 'childList' || mutation.type === 'characterData') {
          // 文本或子节点发生变化
          // console.log('内容发生变化:', selectedValue.innerHTML, selectedValue.dataset.fid, searchValue.value);
          updatePreviewByFetch(selectedValue.closest('.step').querySelector('.step-preview'), selectedValue.dataset.fid, searchValue.value);
        }
      });
    });

    // 下拉树选中框事件监听
    observer.observe(selectedValue, {
      childList: true,       // 监听子节点变化（添加、删除）
      characterData: true,   // 监听文本节点变化
      subtree: true          // 监听所有后代节点
    });

    // 搜索文本框事件监听
    searchValue.addEventListener('input', (event) => {
      // 清除之前的定时器
      clearTimeout(searchValue.dataset.timeoutid);

      // 设置新的定时器
      searchValue.dataset.timeoutid = setTimeout(() => {
        // 定时器到期后执行的处理逻辑
        // console.log('内容发生变化:', selectedValue.innerHTML, selectedValue.dataset.fid, searchValue.value);
        updatePreviewByFetch(selectedValue.closest('.step').querySelector('.step-preview'), selectedValue.dataset.fid, searchValue.value);
      }, 1000);
    });

    // 控制下拉列表的展开/收起
    treeSelect.addEventListener('click', function () {
      treeSelect.classList.toggle('open');
      // open/close 下拉列表
      // const options = treeSelect.querySelector('.tree-options');
      // let lis = options.querySelectorAll('li');
      // let selectValue = treeSelect.querySelector('.selected-value').textContent;
      // if (selectValue === '' || selectValue === '根目录' || selectValue === '/') {
      //   return;
      // }

      // let paths = selectValue.split('/');
      // paths = paths.filter(Boolean);
      // if (paths.length < 1) {
      //   return;
      // }

      // while (paths.length > 1) {
      //   for (let i = 0; i < lis.length; i++) {
      //     if (lis[i].textContent === paths[0]) {
      //       if (!lis[i].classList.contains('has-children')) {
      //         return;
      //       }

      //       lis[i].classList.add('open');
      //       lis = lis[i].querySelectorAll('li');
      //       paths.shift();
      //       break;
      //     }
      //   }
      // }
    });

    // 点击页面其他区域关闭下拉列表
    document.addEventListener('click', function (event) {
      if (!treeSelect.contains(event.target)) {
        treeSelect.classList.remove('open');
      }
    });

    // selectedValue.dataset.fid = '0';
    // selectedValue.dataset.tooltip = '/';
    // selectedValue.textContent = '根目录';

    const treeOptions = treeSelect.querySelector('.tree-options');
    await createTree(null, treeOptions);
  }

  // 请求后端更新下拉树
  async function fetchChildren(parentId) {
    if (parentId === null) {
      parentId = "0";
    }
    return await fetch('/list/p123', { method: 'POST', body: JSON.stringify({ cid: parentId.toString() }) }).then(res => res.json());
    // return new Promise((resolve, reject) => {
    //   // 模拟延迟，增加真实感
    //   setTimeout(() => {
    //     let data;
    //     if (parentId === null) {
    //       data = [
    //         { id: 1, name: '第一项', hasChildren: true },
    //         { id: 2, name: '第二项', hasChildren: false },
    //         { id: 3, name: '第三项', hasChildren: true }
    //       ]
    //     } else if (parentId === 1) {
    //       data = [
    //         { id: 11, name: '第一项-子项1', hasChildren: false },
    //         { id: 12, name: '第一项-子项2', hasChildren: true }
    //       ]
    //     } else if (parentId === 3) {
    //       data = [
    //         { id: 31, name: '第三项-子项1', hasChildren: false },
    //         { id: 32, name: '第三项-子项2', hasChildren: false }
    //       ]
    //     } else if (parentId === 12) {
    //       data = [
    //         { id: 121, name: '第一项-子项2-子项1', hasChildren: false },
    //       ]
    //     } else {
    //       data = [];
    //     }
    //     resolve(data); // 模拟响应数据
    //   }, 500); // 模拟 500ms 延迟
    // });
  }

  // 递归生成树形结构
  async function createTree(parentId, ulElement, path = []) {
    const children = await fetchChildren(parentId);
    if (children === null || children.length === 0) {
      return;
    }

    children.forEach(item => {
      const li = document.createElement('li');
      li.textContent = item.filename;
      li.dataset.id = item.fileId;

      // 点击事件处理函数
      li.addEventListener('click', async function (event) {
        event.stopPropagation(); // 阻止事件冒泡到父级 li 和 treeSelect
        const treeSelect = document.querySelector('.tree-select');
        const selectedValue = treeSelect.querySelector('.selected-value');
        // selectedValue.textContent = item.filename; // 更新 selectedValue

        const currentPath = [...path, item.filename]
        let selectedPath = currentPath.join('/');
        if (!selectedPath.startsWith('/')) {
          selectedPath = '/' + selectedPath;
        }
        // 更新 tooltip 显示完整路径
        selectedValue.dataset.tooltip = selectedPath;
        selectedValue.dataset.fid = item.fileId.toString();
        selectedValue.textContent = selectedPath;

        if (item.type === 1) {
          if (!li.classList.contains('open')) {
            const subTreeUl = document.createElement('ul');
            subTreeUl.classList.add('sub-tree');
            await createTree(item.fileId, subTreeUl, currentPath);
            li.appendChild(subTreeUl);
            li.classList.add('open');
          } else {
            li.classList.remove('open');
            li.querySelector('.sub-tree').remove()
          }
        } else {
          treeSelect.classList.remove('open');
        }
      });

      if (item.type === 1) {
        li.classList.add("has-children");
      }
      ulElement.appendChild(li);
    })
  }

  function setupStepEvents(elem) {
    const nextBtn = elem.querySelector('.next-btn');
    const prevBtn = elem.querySelector('.prev-btn');
    const switchCbx = elem.querySelector('.switch-checkbox');
    const removeBtn = elem.querySelector('.remove-step-btn');
    const submitBtn = elem.querySelector('.submit-btn');
    // const modeSel = elem.querySelector('.mode-sel');
    const viewBtn = elem.querySelector('.preview-btn');

    viewBtn.addEventListener('click', function () {
      const step = this.closest('.step');
      if (stepIndex(step.dataset.stepid) <= 0) {
        updatePreviewByFetch(step.querySelector('.step-preview'), step.querySelector('.selected-value').dataset.fid, step.querySelector('.search-value').value);
      } else {
        updatePreviewByInherit(step.querySelector('.step-preview'));
      }
    });

    nextBtn.addEventListener('click', function () {
      const stepElement = this.closest('.step');
      let idx = stepIndex(stepElement.dataset.stepid);
      if (idx < 0) {
        return;
      }

      if (idx === steps.length - 1) {
        next = idx;
      } else {
        next = idx + 1;
        showStep(steps[next].dataset.stepid);
        steps[next].querySelector('.preview-btn').click();
      }
    });

    prevBtn.addEventListener('click', function () {
      const stepElement = this.closest('.step');
      let idx = stepIndex(stepElement.dataset.stepid);
      if (idx < 0) {
        return;
      }

      if (idx === 0) {
        prev = idx;
      } else {
        prev = idx - 1;
        showStep(steps[prev].dataset.stepid);
        steps[prev].querySelector('.preview-btn').click();
      }
    });

    if (switchCbx) {
      switchCbx.addEventListener('change', function () {
        const stepElement = this.closest('.step');
        if (this.checked) {
          stepElement.dataset.disabled = false;
        } else {
          stepElement.dataset.disabled = true;
        }
      });
    }

    if (removeBtn) {
      removeBtn.addEventListener('click', function () {
        const stepElement = this.closest('.step');
        removeStep(stepElement.dataset.stepid);
      });
    }

    if (submitBtn) {
      submitBtn.addEventListener('click', function () {
        const step = this.closest('.step');
        const previewContainer = step.querySelector('.preview-container');
        const fileitems = previewContainer.querySelectorAll('.file-item');
        let data = [];
        fileitems.forEach(item => {
          if (item.dataset.selected == "1") {
            data.push({
              fid: item.dataset.fid,
              name: item.dataset.nfname,
              type: item.dataset.ftype,
              path: item.dataset.fpath,
            });
          }
        });

        fetch('/rename/p123', { method: 'POST', body: JSON.stringify(data) }).then(response => response.json()).then(data => {
          console.log(data);
          let succ = data.filter(item => item.status === "1").length;
          let fail = data.filter(item => item.status === "0").length;
          console.log(`重命名成功：${succ} 条，失败：${fail} 条`);
          // 弹出操作结果
          showModal('success', `操作成功：${succ} 条，失败：${fail} 条`);
          // 回到第一个步骤刷新最新结果
          showStep(steps[0].dataset.stepid);
          // 刷新页面
        }).catch(error => {
          console.error('Error:', error);
          showModal('error', '提交失败:' + error.message);
        });
      });
    }
  }

  function stepIndex(stepid = "") {
    if (steps.length === 0) {
      return -1;
    }
    if (stepid === "") {
      stepid = steps.find(item => item.classList.contains('active')).dataset.stepid;
    }
    return steps.indexOf(steps.find(item => item.dataset.stepid === stepid));
  }

  // 更新步骤信息
  function updateStepInfo() {
    stepInfoDiv.textContent = `当前步骤: ${stepIndex() + 1} / 总步骤: ${steps.length}`;
  }


  function showStep(stepid) {
    let cstep = steps[0];
    steps.forEach(step => {
      step.classList.remove('active');
      if (step.dataset.stepid == stepid) {
        step.classList.add('active');
        cstep = step;
      }
    });
    updateStepInfo(); // 显示步骤信息

    // 更新步骤预览
    if (cstep.dataset.type.startsWith("default")) {
      updatePreviewByFetch(cstep.querySelector('.step-preview'), cstep.querySelector('.selected-value').dataset.fid, cstep.querySelector('.search-value').value);
    } else {
      updatePreviewByInherit(cstep.querySelector('.step-preview'));
    }

  }

  function showModal(type, message) {
    const overlay = document.querySelector('.modal-overlay');
    const modal = overlay.querySelector('.modal');
    const text = modal.querySelector('.modal-text');

    // 清除之前的类名
    modal.classList.remove('success', 'error', 'info');
    // 根据传入类型添加对应的类名
    if (type === 'success') {
      modal.classList.add('success')
    } else if (type === 'error') {
      modal.classList.add('error')
    } else {
      modal.classList.add('info')
    }

    text.textContent = message;
    overlay.classList.add('active');
    // 设置定时器自动隐藏
    setTimeout(hideModal, 3000);
  }


  function hideModal() {
    const overlay = document.querySelector('.modal-overlay');
    overlay.classList.remove('active');
  }

  // let strs = ["月球陨落.Moonfall.2021.2160p.WEB-DL.x26665.10bit.HDR.DDP5.1.Atmos-NOGRP.mkv"];
  // let action = new Action("add", "regexp", "2\\d", "ABC", 0);
  // let nstrs = action.rename(strs);
  // console.log(strs);
  // console.log(nstrs);
  // return;

  // 初始化第一个步骤
  addStep();
  initializeTree();


  // 下拉选框点击事件
  dropdownContent.querySelectorAll('a').forEach(item => {
    item.addEventListener('click', function (e) {
      e.preventDefault();
      addStep(this.dataset.type);
    })
  });
});

function getAllIndex(str, subStr) {
  const indexes = [];
  let index = str.indexOf(subStr);
  while (index !== -1) {
    indexes.push(index);
    index = str.indexOf(subStr, index + 1); // 从下一个位置继续查找
  }
  return indexes;
}

function getAllRegIndex(str, regexp, num = -1) {
  const indexes = [];
  const reg = new RegExp(regexp, 'g');
  const matches = str.matchAll(reg);

  for (const match of matches) {
    if (num > 0 && indexes.length >= num) {
      break;
    }
    indexes.push(match.index);
  }
  return indexes;
}

class Action {
  constructor(kind, mode, expr, value, num) {
    this.kind = kind;
    this.mode = mode;
    this.expr = expr;
    this.value = value;
    this.num = num;
  }

  filter(strs) {
    switch (this.kind) {
      case "include":
        return this.include(strs);
      case "exclude":
        return this.exclude(strs);
      default:
        return strs;
    }
  }

  rename(strs) {
    switch (this.kind) {
      case "add":
        return this.add(strs);
      case "del":
        return this.rename(strs);
      default:
        return strs;
    }
  }

  include(strs) {
    let output = [];
    let num = this.num;

    if (this.num <= 0) {
      num = -1;
    }

    for (let idx = 0; idx < strs.length; idx++) {
      const str = strs[idx];

      if (idx >= num && num !== -1) {
        output.push(str);
        continue; // 跳过本次循环，进行下次循环
      }

      switch (this.mode) {
        case "equal":
          if (this.expr === str) {
            output.push(str);
          }
          break;
        case "contain":
          if (str.includes(this.expr)) {
            output.push(str);
          }
          break;
        case "prefix":
          if (str.startsWith(this.expr)) {
            output.push(str);
          }
          break;
        case "suffix":
          if (str.endsWith(this.expr)) {
            output.push(str);
          }
          break;
        case "regexp":
          try {
            const regex = new RegExp(this.expr); // 尝试创建正则表达式
            if (regex.test(str)) {
              output.push(str);
            }
          } catch (err) {
            output.push(str); // 如果正则表达式无效，将 str 添加到输出
          }
          break;
        default:
          output.push(str);
      }
    }

    return output;
  }

  exclude(strs) {
    let output = [];
    let num = this.num;

    if (this.num <= 0) {
      num = -1;
    }

    for (let idx = 0; idx < strs.length; idx++) {
      const str = strs[idx];

      if (idx >= num && num !== -1) {
        output.push(str);
        continue; // 跳过本次循环，进行下次循环
      }

      switch (this.mode) {
        case "equal":
          if (this.expr != str) {
            output.push(str);
          }
          break;
        case "contain":
          if (!str.includes(this.expr)) {
            output.push(str);
          }
          break;
        case "prefix":
          if (!str.startsWith(this.expr)) {
            output.push(str);
          }
          break;
        case "suffix":
          if (!str.endsWith(this.expr)) {
            output.push(str);
          }
          break;
        case "regexp":
          try {
            const regex = new RegExp(this.expr); // 尝试创建正则表达式
            if (!regex.test(str)) {
              output.push(str);
            }
          } catch (err) {
            output.push(str); // 如果正则表达式无效，将 str 添加到输出
          }
          break;
        default:
          output.push(str);
      }
    }
    return output;
  }

  add(strs) {
    let output = [];
    for (let str of strs) {
      switch (this.mode) {
        case "plain":
          var num = this.num;
          if (this.num <= 0) {
            num = -1;
          }
          var idxs = getAllIndex(str, this.expr);
          for (let i = 0; i < idxs.length; i++) {
            const idx = idxs[i];
            if (i < num || num == -1) {
              str = str.slice(0, idx + i * this.value.length) + this.value + str.slice(idx + i * this.value.length);
            }
          }
          output.push(str);
          break;
        case "prefix":
          output.push(this.value + str);
          break;
        case "suffix":
          output.push(str + this.value);
          break;
        case "index":
          var start = 1;
          //const istrs = [...str]; // 将字符串转换为 Unicode 字符数组
          const ranger = this.expr.split(":");
          if (ranger[0] !== "") {
            start = parseInt(ranger[0]);
            if (isNaN(start)) {
              output.push(str)
              break
            }
          }

          if (start == 0) {
            start = 1;
          }
          if (start < 0) {
            start = str.length + 1 + start;
          }

          if (start > str.length) {
            output.push(str + this.value);
            break
          }
          if (str.length > 0) {
            start -= 1;
          }

          output.push(str.slice(0, start) + this.value + str.slice(start));
          break;
        case "regexp":
          var num = this.num;
          if (this.num <= 0) {
            num = -1;
          }

          const regex = new RegExp(this.expr);
          var idxs = getAllRegIndex(str, regex, num);
          for (let i = 0; i < idxs.length; i++) {
            const idx = idxs[i]
            str = str.slice(0, idx + i * this.value.length) + this.value + str.slice(idx + i * this.value.length);
          }
          output.push(str);
          break;
        default:
          output.push(str);
      }
    }
    return output;
  }
}