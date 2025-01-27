document.addEventListener('DOMContentLoaded', function () {
  const steps = document.querySelectorAll('.step');
  const nextBtns = document.querySelectorAll('.next-btn');
  const prevBtns = document.querySelectorAll('.prev-btn');
  const submitBtn = document.querySelector('.submit-btn');
  let currentStep = 1;

  function showStep(stepNumber) {
    steps.forEach(step => {
      step.classList.remove('active');
      if (step.dataset.step == stepNumber) {
        step.classList.add('active');
      }
    });
  }

  showStep(currentStep);

  nextBtns.forEach(btn => {
    btn.addEventListener('click', function () {
      if (currentStep < steps.length) {
        currentStep++;
        showStep(currentStep)
      }
    });
  });

  prevBtns.forEach(btn => {
    btn.addEventListener('click', function () {
      if (currentStep > 1) {
        currentStep--;
        showStep(currentStep);
      }
    });
  });

  if (submitBtn) {
    submitBtn.addEventListener('click', function () {
      // 获取步骤1的信息
      const name = document.getElementById('name').value;
      const email = document.getElementById('email').value;

      // 获取步骤2的信息
      const favColor = document.getElementById('fav-color').value;
      const pet = document.querySelector('input[name="pet"]:checked') ? document.querySelector('input[name="pet"]:checked').value : "未选择";
      // 渲染到最后一步
      const confirmationArea = document.getElementById('confirmation-area');
      confirmationArea.innerHTML = `
              <p>姓名: ${name}</p>
              <p>邮箱: ${email}</p>
              <p>喜欢的颜色: ${favColor}</p>
              <p>喜欢宠物: ${pet}</p>
          `;
      alert('信息已提交!');
    });
  }


});